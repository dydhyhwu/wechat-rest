package whapp

import (
	"github.com/opentdp/wrest-chat/wclient/whapp/gitea"
	"github.com/opentdp/wrest-chat/wclient/whapp/github"
	"github.com/opentdp/wrest-chat/wclient/whapp/text"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Handler(header http.Header, app string, msg string) string {

	var res string
	var err error

	switch app {
	case "github":
		res, err = github.HandleWebhook(header, msg)
	case "gitea":
		res, err = gitea.HandleWebhook(header, msg)
	case "text":
		res, err = text.HandleWebhook(msg)
	default:
		log.Warnf("unsupport app: %s", app)
	}

	if err != nil {
		return err.Error()
	}

	return res
}
