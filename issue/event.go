package issue

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
)

type IssueCommentEvent github.IssueCommentEvent

// if command is /pay fanux 10 ,type is : pay Command is : fanux 10
type Command struct {
	Type    string // like pay
	Command string // like 10, full command not contain type
}

type IssueEvent struct {
	*IssueCommentEvent
	Command *Command
	Client  *github.Client
}

var robot map[string]Robot

type Robot interface {
	Process(event IssueEvent) error
}

// github user config
type Config struct {
	UserName string
	Password string
	Token    string
}

func NewConfig(user string, passwd string) Config {
	if user == "" {
		user = os.Getenv("GITHUB_USER")
	}
	if passwd == "" {
		passwd = os.Getenv("GITHUB_PASSWD")
	}
	return Config{UserName: user, Password: passwd}
}

func Process(config Config, event IssueCommentEvent) error {
	tp := github.BasicAuthTransport{
		Username: config.UserName,
		Password: config.Password,
	}
	client := github.NewClient(tp.Client())
	var body *string
	if event.Comment == nil {
		body = event.Issue.Body
	} else {
		body = event.Comment.Body
	}
	//decode commands
	commands := decodeFromBody(body)
	fmt.Println("commands from body:", commands)

	for _, command := range commands {
		issueEvent := IssueEvent{
			&event,
			command,
			client,
		}
		fmt.Println("process command", command.Type, command.Command)
		if v, ok := robot[command.Type]; ok {
			err := v.Process(issueEvent)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Regist user robot
func Regist(command string, r Robot) {
	if robot == nil {
		robot = make(map[string]Robot)
	}
	robot[command] = r
}

func decodeFromBody(body *string) []*Command {
	var res []*Command
	lines := strings.Split(*body, "\r\n")
	for _, line := range lines {
		if !validCommand(line) {
			continue
		}
		res = append(res, decodeCommand(line))
	}
	return res
}

func validCommand(s string) bool {
	for _, b := range s {
		t := byte(b)
		if t != ' ' && t != '/' {
			return false
		}
		if t == '/' {
			return true
		}
		if t == ' ' {
			continue
		}
	}
	return false
}

// decode /pay 10 like command
func decodeCommand(s string) *Command {
	command := &Command{}
	var i, j int
	fmt.Printf("decode cmd: %s\n", s)
	for i = range s {
		if byte(s[i]) == '/' {
			break
		}
	}
	var flag bool
	for j = i; j < len(s); j++ {
		if !flag && (byte(s[j]) == ' ' || j==len(s)-1){
			flag = true
			if s[j] == ' ' {
				command.Type = s[i:j]
			} else {
				command.Type = s[i:j+1]
				break
			}
		}
		if flag && byte(s[j]) != ' ' {
			command.Command = s[j:]
			break
		}
	}
	return command
}
