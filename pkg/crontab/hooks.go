package crontab

type HookFunc func(jobID string)

type Hooks struct {
	BeforeJob []HookFunc
	AfterJob  []HookFunc
}

func (h *Hooks) RunBefore(jobID string) {
	for _, hook := range h.BeforeJob {
		hook(jobID)
	}
}

func (h *Hooks) RunAfter(jobID string) {
	for _, hook := range h.AfterJob {
		hook(jobID)
	}
}

func (h *Hooks) AddBeforeHook(hook HookFunc) {
	h.BeforeJob = append(h.BeforeJob, hook)
}

func (h *Hooks) AddAfterHook(hook HookFunc) {
	h.AfterJob = append(h.AfterJob, hook)
}
