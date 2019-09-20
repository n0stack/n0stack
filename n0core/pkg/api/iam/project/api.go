package project

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"n0st.ac/n0stack/n0core/pkg/api/iam/authn"
	stdapi "n0st.ac/n0stack/n0core/pkg/api/standard_api"
	"n0st.ac/n0stack/n0core/pkg/datastore"
	grpcutil "n0st.ac/n0stack/n0core/pkg/util/grpc"
	piam "n0st.ac/n0stack/n0proto.go/iam/v1alpha"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type ProjectAPI struct {
	datastore datastore.Datastore
	userapi   piam.UserServiceClient
}

func CreateProjectAPI(ds datastore.Datastore, userapi piam.UserServiceClient) *ProjectAPI {
	return &ProjectAPI{
		datastore: ds.AddPrefix("iam/project"),
		userapi:   userapi,
	}
}

func (a *ProjectAPI) ListProjects(ctx context.Context, req *piam.ListProjectsRequest) (*piam.ListProjectsResponse, error) {
	return ListProjects(ctx, req, a.datastore)
}

func (a *ProjectAPI) GetProject(ctx context.Context, req *piam.GetProjectRequest) (*piam.Project, error) {
	u, _, err := GetProject(ctx, a.datastore, req.Name)
	return u, err
}

func (a *ProjectAPI) CreateProject(ctx context.Context, req *piam.CreateProjectRequest) (*piam.Project, error) {
	if err := stdapi.ValidateName(req.Name); err != nil {
		return nil, err
	}

	if _, _, err := GetProject(ctx, a.datastore, req.Name); err != nil {
		if grpc.Code(err) != codes.NotFound {
			return nil, err
		}
	}

	username, err := authn.GetConnectingAccountName(ctx)
	if err != nil {
		return nil, err
	}

	project := &piam.Project{
		Name:        req.Name,
		Annotations: req.Annotations,
		Labels:      req.Labels,

		Membership: map[string]piam.ProjectMembership{
			username: piam.ProjectMembership_OWNER,
		},
	}

	if _, err := ApplyProject(ctx, a.datastore, project, 0); err != nil {
		return nil, err
	}

	return project, nil
}

func (a *ProjectAPI) UpdateProject(ctx context.Context, req *piam.CreateProjectRequest) (*piam.Project, error) {
	project, version, err := GetProject(ctx, a.datastore, req.Name)
	if err != nil {
		return nil, err
	}

	username, err := authn.GetConnectingAccountName(ctx)
	if err != nil {
		return nil, err
	}

	if project.Membership[username] != piam.ProjectMembership_OWNER {
		return nil, grpc.Errorf(codes.PermissionDenied, "Owner account can only UpdateProject()")
	}

	project = &piam.Project{
		Name:        req.Name,
		Annotations: req.Annotations,
		Labels:      req.Labels,
	}

	if _, err := ApplyProject(ctx, a.datastore, project, version); err != nil {
		return nil, err
	}

	return project, nil
}

func (a *ProjectAPI) DeleteProject(ctx context.Context, req *piam.DeleteProjectRequest) (*empty.Empty, error) {
	project, version, err := GetProject(ctx, a.datastore, req.Name)
	if err != nil {
		if grpc.Code(err) != codes.NotFound {
			return &empty.Empty{}, nil
		}

		return nil, err
	}

	username, err := authn.GetConnectingAccountName(ctx)
	if err != nil {
		return nil, err
	}

	if project.Membership[username] != piam.ProjectMembership_OWNER {
		return nil, grpc.Errorf(codes.PermissionDenied, "Owner account can only DeleteProject()")
	}

	if err := DeleteProject(ctx, a.datastore, project.Name, version); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (a *ProjectAPI) AddProjectMembership(ctx context.Context, req *piam.AddProjectMembershipRequest) (*piam.Project, error) {
	if req.Membership == piam.ProjectMembership_PROJECT_MEMBERSHIP_UNSPECIFIED {
		return nil, stdapi.ValidationError("member_ship", "necessary to specify any membership")
	}

	project, version, err := GetProject(ctx, a.datastore, req.UserName)
	if err != nil {
		return nil, err
	}

	username, err := authn.GetConnectingAccountName(ctx)
	if err != nil {
		return nil, err
	}

	if project.Membership[username] != piam.ProjectMembership_OWNER {
		return nil, grpc.Errorf(codes.PermissionDenied, "Owner account can only AddProjectMembership()")
	}

	_, err = a.userapi.GetUser(ctx, &piam.GetUserRequest{Name: req.UserName})
	if err != nil {
		if grpc.Code(err) == codes.Internal {
			return nil, err
		}

		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	if project.Membership == nil {
		project.Membership = make(map[string]piam.ProjectMembership)
	}
	project.Membership[req.UserName] = req.Membership

	if _, err := ApplyProject(ctx, a.datastore, project, version); err != nil {
		return nil, err
	}

	return project, nil
}
func (a *ProjectAPI) DeleteProjectMembership(ctx context.Context, req *piam.DeleteProjectMembershipRequest) (*piam.Project, error) {
	project, version, err := GetProject(ctx, a.datastore, req.UserName)
	if err != nil {
		return nil, err
	}

	username, err := authn.GetConnectingAccountName(ctx)
	if err != nil {
		return nil, err
	}

	if project.Membership[username] != piam.ProjectMembership_OWNER {
		return nil, grpc.Errorf(codes.PermissionDenied, "Owner account can only DeleteProjectMembership()")
	}

	if project.Membership == nil {
		return nil, grpcutil.Errorf(codes.NotFound, "publicKey '%s' does not exist", req.UserName)
	}
	if _, ok := project.Membership[req.UserName]; !ok {
		return nil, grpcutil.Errorf(codes.NotFound, "publicKey '%s' does not exist", req.UserName)
	}

	delete(project.Membership, req.UserName)

	ownerExists := false
	for _, m := range project.Membership {
		if m == piam.ProjectMembership_OWNER {
			ownerExists = true
		}
	}
	if !ownerExists {
		return nil, grpcutil.Errorf(codes.OutOfRange, "Project owner are gone")
	}

	if _, err := ApplyProject(ctx, a.datastore, project, version); err != nil {
		return nil, err
	}

	return project, nil
}
