package iam

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	piam "n0st.ac/n0stack/iam/v1alpha"
	stdapi "n0st.ac/n0stack/n0core/pkg/api/stdapi"
	"n0st.ac/n0stack/n0core/pkg/datastore"
	"n0st.ac/n0stack/n0core/pkg/driver/n0stack/auth"
	grpcutil "n0st.ac/n0stack/n0core/pkg/util/grpc"
	structutil "n0st.ac/n0stack/n0core/pkg/util/struct"
)

type ProjectAPI struct {
	datastore datastore.Datastore
	userapi   piam.UserServiceClient

	auth *auth.AuthenticationServiceProvier
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
	if err := stdapi.ValidateName(req.Project.Name); err != nil {
		return nil, err
	}

	if _, _, err := GetProject(ctx, a.datastore, req.Project.Name); err != nil {
		if grpc.Code(err) != codes.NotFound {
			return nil, err
		}
	}

	username, err := stdapi.GetAuthenticatedUserName(ctx, a.auth)
	if err != nil {
		return nil, err
	}

	project := &piam.Project{
		Name:        req.Project.Name,
		Annotations: req.Project.Annotations,
		Labels:      req.Project.Labels,

		Membership: map[string]piam.ProjectMembership{
			username: piam.ProjectMembership_OWNER,
		},
	}

	if _, err := ApplyProject(ctx, a.datastore, project, 0); err != nil {
		return nil, err
	}

	return project, nil
}

func (a *ProjectAPI) UpdateProject(ctx context.Context, req *piam.UpdateProjectRequest) (*piam.Project, error) {
	username, err := stdapi.GetAuthenticatedUserName(ctx, a.auth)
	if err != nil {
		return nil, err
	}

	project, version, err := GetProject(ctx, a.datastore, req.Project.Name)
	if err != nil {
		return nil, err
	}

	if err := stdapi.IsOwner(project, username); err != nil {
		return nil, err
	}

	if err := structutil.UpdateWithMaskUsingJson(project, req.Project, req.UpdateMask.Paths); err != nil {
		return nil, stdapi.UpdateMaskError(err)
	}

	if _, err := ApplyProject(ctx, a.datastore, project, version); err != nil {
		return nil, err
	}

	return project, nil
}

func (a *ProjectAPI) DeleteProject(ctx context.Context, req *piam.DeleteProjectRequest) (*empty.Empty, error) {
	username, err := stdapi.GetAuthenticatedUserName(ctx, a.auth)
	if err != nil {
		return nil, err
	}

	project, version, err := GetProject(ctx, a.datastore, req.Name)
	if err != nil {
		if grpc.Code(err) != codes.NotFound {
			return &empty.Empty{}, nil
		}

		return nil, err
	}

	if err := stdapi.IsOwner(project, username); err != nil {
		return nil, err
	}

	if err := DeleteProject(ctx, a.datastore, project.Name, version); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (a *ProjectAPI) AddProjectMembership(ctx context.Context, req *piam.AddProjectMembershipRequest) (*piam.Project, error) {
	if req.Membership == piam.ProjectMembership_PROJECT_MEMBERSHIP_UNSPECIFIED {
		return nil, stdapi.ValidationError("membership", "necessary to specify any membership")
	}

	username, err := stdapi.GetAuthenticatedUserName(ctx, a.auth)
	if err != nil {
		return nil, err
	}

	project, version, err := GetProject(ctx, a.datastore, req.UserName)
	if err != nil {
		return nil, err
	}

	if err := stdapi.IsOwner(project, username); err != nil {
		return nil, err
	}

	_, err = a.userapi.GetUser(ctx, &piam.GetUserRequest{Name: req.UserName})
	if err != nil {
		return nil, grpcutil.GrpcWrapf(err, "failed to check user")
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

	username, err := stdapi.GetAuthenticatedUserName(ctx, a.auth)
	if err != nil {
		return nil, err
	}

	if err := stdapi.IsOwner(project, username); err != nil {
		return nil, err
	}

	if err := stdapi.IsMember(project, req.UserName); err != nil {
		return nil, err
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
