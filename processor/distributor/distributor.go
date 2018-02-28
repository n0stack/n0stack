package distributor

// import (
// 	"github.com/n0stack/n0core/gateway"
// 	"github.com/n0stack/n0core/message"
// 	"github.com/n0stack/n0core/model"
// 	"github.com/n0stack/n0core/repository"
// )

// // Distributor is a processor which schedule resource from spec message
// // and distribute the notification messages to agents.
// //
// // 1. Receive a message from gateway.
// // 2. Schdule resource from spec message.
// // 3. Send notification message to the agent on scheduled host with gateway.
// // 4. Send notification message to aggreagater to notify result.
// //
// // Args:
// // 	repository: Data store to schedule resource.
// // 	notification: Gateway to notify result to agent and aggregater.
// //
// // Example:
// type Distributor struct {
// 	Repository repository.Repository
// 	Notifier   gateway.Gateway
// }

// func (d Distributor) ProcessMessage(m message.AbstractMessage) {
// 	s, ok := m.(message.Spec)
// 	if !ok {
// 		return
// 	}

// 	for _, am := range s.Models {
// 		mo := am.ToModel()

// 		// check whether model is already applied or not
// 		if d.modelIsApllied(mo) {
// 			continue
// 		}

// 		if _, err
// 	}
// }

// func (d Distributor) modelIsApllied(m *model.Model) bool {
// 	_, err := d.Repository.DigModel(m.ID, "APPLIED", 0)
// 	if err == nil {
// 		return true
// 	}

// 	return false
// }
