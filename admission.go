package testac

import (
	"errors"
	"io"
        "strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/kubernetes/pkg/api"
)

func init() {
	admission.RegisterPlugin("Testac", func(config io.Reader) (admission.Interface, error) {
		return NewTestac(), nil
	})
}

type testac struct {
	*admission.Handler
}

func (a *testac) Admit(attributes admission.Attributes) (err error) {
	if len(attributes.GetSubresource()) != 0 || attributes.GetResource().GroupResource() != api.Resource("pods") {
		return nil
	}
	pod, ok := attributes.GetObject().(*api.Pod)
	if !ok {
		return apierrors.NewBadRequest("Issue casting to pod.")
	}

        name := pod.GenerateName

        if strings.HasPrefix(name, "toocool") || strings.HasPrefix(attributes.GetName(), "toocool") {
            return admission.NewForbidden(attributes, errors.New("Starts with toocool."))
        }

	return nil
}

func NewTestac() admission.Interface {
	return &testac{
		Handler: admission.NewHandler(admission.Create, admission.Update),
	}
}
