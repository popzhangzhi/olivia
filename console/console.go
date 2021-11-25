package console

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"

	"github.com/olivia-ai/olivia/analysis"
	"github.com/olivia-ai/olivia/locales"
	"github.com/olivia-ai/olivia/network"
	"github.com/olivia-ai/olivia/server"
	"github.com/olivia-ai/olivia/user"
	"github.com/olivia-ai/olivia/util"
	"github.com/pkg/errors"
)

type Service interface {
	Serve(neuralNetworks map[string]network.Network)
	Reply(locale, content string, neuralNetworks map[string]network.Network) ([]byte, error)
	RandomToken() (string, error)
}

var _ Service = (*ServiceImpl)(nil)

type ServiceImpl struct {
}

func NewService() *ServiceImpl {
	return &ServiceImpl{}
}

// Serve 运行cli模式交互
func (s ServiceImpl) Serve(neuralNetworks map[string]network.Network) {

	ch := make(chan string)
	go func() {
		for {
			text, ok := <-ch
			if !ok {
				util.CliError("read channel error")
			}
			response, err := s.Reply("en", text, neuralNetworks)
			if err != nil {
				util.CliError(err.Error())
				return
			}
			util.CliInfo(string(response))
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			util.CliError(err.Error())
		}
		ch <- scanner.Text()
	}

}

// Reply takes the entry message and returns an array of bytes for the answer
func (s ServiceImpl) Reply(locale, content string, neuralNetworks map[string]network.Network) ([]byte, error) {
	var (
		responseSentence, responseTag string
		err                           error
	)

	token, err := s.RandomToken()
	if err != nil {
		return nil, errors.Wrap(err, "random token error")
	}
	// Send a message from res/datasets/messages.json if it is too long
	if len(content) > 500 {
		responseTag = "too long"
		responseSentence = util.GetMessage(locale, responseTag)
	} else {
		// If the given locale is not supported yet, set english
		if !locales.Exists(locale) {
			locale = "en"
		}

		responseTag, responseSentence = analysis.NewSentence(
			locale, content,
		).Calculate(*server.Cache, neuralNetworks[locale], token)
	}

	// Marshall the response in json
	response := server.ResponseMessage{
		Content:     responseSentence,
		Tag:         responseTag,
		Information: user.GetUserInformation(token),
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	return bytes, nil
}

// RandomToken 随机30byte token
func (s ServiceImpl) RandomToken() (string, error) {
	b := make([]byte, 30)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
