/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	playv1alpha1 "github.com/roehrich-hpe/airplane-sim/api/v1alpha1"
)

var _ = Describe("Rudder Unit Tests", func() {

	var (
		key    types.NamespacedName
		rudder *playv1alpha1.Rudder
	)

	BeforeEach(func() {
		key = types.NamespacedName{
			Name:      "rudder-" + uuid.New().String()[0:8],
			Namespace: corev1.NamespaceDefault,
		}

		rudder = &playv1alpha1.Rudder{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
			},
			Spec: playv1alpha1.RudderSpec{
				//Position: "left",
			},
		}

		Expect(k8sClient.Create(context.TODO(), rudder)).To(Succeed())
	})

	AfterEach(func() {
		Expect(k8sClient.Delete(context.TODO(), rudder)).To(Succeed())

		// Test env limitation: Wait until the cached object is gone.
		expected := &playv1alpha1.Rudder{}
		Eventually(func() error {
			return k8sClient.Get(context.TODO(), key, expected)
		}).ShouldNot(Succeed())
	})

	It("Moves rudder to new position", func() {
		Eventually(func(g Gomega) {
			g.Expect(k8sClient.Get(context.TODO(), key, rudder)).To(Succeed())
			g.Expect(rudder.Status.Position).To(Equal(rudder.Spec.Position))
		}).Should(Succeed())
	})

	It("Moves to rudder to next position", func() {
		Expect(k8sClient.Get(context.TODO(), key, rudder)).To(Succeed())

		var wanted string
		switch rudder.Spec.Position {
		case "neutral":
			wanted = "left"
		case "left":
			wanted = "right"
		case "right":
			wanted = "neutral"

		}
		Eventually(func(g Gomega) {
			g.Expect(k8sClient.Get(context.TODO(), key, rudder)).To(Succeed())
			rudder.Spec.Position = wanted
			g.Expect(k8sClient.Update(context.TODO(), rudder)).To(Succeed())
		}).Should(Succeed())

		expected := &playv1alpha1.Rudder{}
		Eventually(func(g Gomega) {
			g.Expect(k8sClient.Get(context.TODO(), key, expected)).To(Succeed())
			g.Expect(expected.Status.Position).To(Equal(wanted))
		}).Should(Succeed())
	})
})
