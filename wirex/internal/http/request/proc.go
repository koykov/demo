package request

type ProcFn func(req *Inner) error

type ProcFnContainer []ProcFn
