package beater

import (
	"time"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
)

const selector = "mesosbeat"
const selectorDetail = "json"

type Mesosbeat struct {
	period 		time.Duration
	urls   		[]string
	rootdir 	string

	MbConfig ConfigSettings
	events   publisher.Client

	done 			chan struct{}
}

func New() *Mesosbeat {
	return &Mesosbeat{
		done: make(chan struct{}),
	}
}

func (mb *Mesosbeat) Config(b *beat.Beat) error {

	err := cfgfile.Read(&mb.MbConfig, "")
	if err != nil {
		logp.Err("Error reading configuration file: %v", err)
		return err
	}

	if mb.MbConfig.Input.Period != nil {
		mb.period = time.Duration(*mb.MbConfig.Input.Period) * time.Second
	} else {
		mb.period = 10 * time.Second
	}

	//define default URL if none provided
	var urlConfig []string
	if mb.MbConfig.Input.URLs != nil {
		urlConfig = mb.MbConfig.Input.URLs
	} else {
		urlConfig = []string{"http://127.0.0.1:5051/metrics/snapshot"}
	}

	if mb.MbConfig.Input.RootDir != "" {
		mb.rootdir = mb.MbConfig.Input.RootDir
	} else {
		mb.rootdir = DefaultRootDir
	}

	mb.urls = make([]string, len(urlConfig))
	for i := 0; i < len(urlConfig); i++ {
		u := urlConfig[i]
		mb.urls[i] = u
	}

	logp.Debug(selector, "Init mesosbeat")
	logp.Debug(selector, "Period %v\n", mb.period)
	logp.Debug(selector, "Watch %v", mb.urls)

	return nil
}

func (mb *Mesosbeat) Setup(b *beat.Beat) error {
	mb.events = b.Events
	mb.done = make(chan struct{})
	return nil
}

func (mb *Mesosbeat) Run(b *beat.Beat) error {
	logp.Debug(selector, "Run mesosbeat")

	//for each url
	for _, u := range mb.urls {
		go func(u string) {
			ticker := time.NewTicker(mb.period)
			defer ticker.Stop()

			for {
				select {
				case <-mb.done:
					goto GotoFinish
				case <-ticker.C:
				}

				timerStart := time.Now()

				logp.Debug(selector, "Fetching agent attributes from: %v", mb.rootdir)
				attributes, err := mb.GetAgentAttributes(mb.rootdir)
				if err != nil {
					logp.Err("Error reading agent attributes: %v", err)
				} else {
					logp.Debug(selectorDetail, "Agent attributes: %+v", attributes)

//					event := common.MapStr{
//						"@timestamp":    common.Time(time.Now()),
//						"type":          "agent_attributes",
//						"attributes": 	attributes,
//					}
//					mb.events.PublishEvent(event)
				}

				logp.Debug(selector, "Fetching agent statistics from: %v", u)
				agent_stats, err := mb.GetAgentStatistics(u)
				if err != nil {
					logp.Err("Error reading agent statistics: %v", err)
				} else {
					logp.Debug(selectorDetail, "Agent statistics: %+v", agent_stats)

					event := common.MapStr{
						"@timestamp":    common.Time(time.Now()),
						"type":          "mesosbeat",
						"attributes":	attributes,
						"statistics":	agent_stats,
					}
					mb.events.PublishEvent(event)
				}

				timerEnd := time.Now()
				duration := timerEnd.Sub(timerStart)
				if duration.Nanoseconds() > mb.period.Nanoseconds() {
					logp.Warn("Ignoring tick(s) due to processing taking longer than one period")
				}
			}

		GotoFinish:
		}(u)
	}

	<-mb.done
	return nil
}

func (mb *Mesosbeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (mb *Mesosbeat) Stop() {
	logp.Debug(selector, "Stop mesosbeat")
	close(mb.done)
}
