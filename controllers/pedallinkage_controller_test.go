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

var _ = Describe("PedalLinkage Unit Tests", func() {

	var (
		key    types.NamespacedName
		pedals *playv1alpha1.Pedals

		rudder *playv1alpha1.Rudder
	)

	BeforeEach(func() {
		// Peddles and rudder have the same value for the name.
		key = types.NamespacedName{
			Name:      uuid.New().String()[0:8],
			Namespace: corev1.NamespaceDefault,
		}

		pedals = &playv1alpha1.Pedals{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
			},
			Spec: playv1alpha1.PedalsSpec{
				//Pressed: "left",
			},
		}

		Expect(k8sClient.Create(context.TODO(), pedals)).To(Succeed())

		rudder = &playv1alpha1.Rudder{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
			},
			Spec: playv1alpha1.RudderSpec{
				//Position: "neutral",
			},
		}

		Expect(k8sClient.Create(context.TODO(), rudder)).To(Succeed())
	})

	AfterEach(func() {
		Expect(k8sClient.Delete(context.TODO(), pedals)).To(Succeed())

		// Test env limitation: Wait until the cached object is gone.
		expected := &playv1alpha1.Pedals{}
		Eventually(func() error {
			return k8sClient.Get(context.TODO(), key, expected)
		}).ShouldNot(Succeed())
	})

	It("Moves pedals to new position", func() {
		Eventually(func(g Gomega) {
			g.Expect(k8sClient.Get(context.TODO(), key, pedals)).To(Succeed())
			g.Expect(pedals.Status.LinkagePosition).To(Equal("neutral"))
		}).Should(Succeed())

		Eventually(func(g Gomega) {
			g.Expect(k8sClient.Get(context.TODO(), key, rudder)).To(Succeed())
			g.Expect(rudder.Status.Position).To(Equal("neutral"))
		}).Should(Succeed())
	})

	It("Moves to pedals to next position", func() {
		Expect(k8sClient.Get(context.TODO(), key, pedals)).To(Succeed())

		var wanted string
		var linkageWanted string
		switch pedals.Spec.Pressed {
		case "none":
			wanted = "left"
			linkageWanted = "left"
		case "left":
			wanted = "right"
			linkageWanted = "right"
		case "right":
			wanted = "none"
			linkageWanted = "neutral"

		}

		By("moving the peddles")
		Eventually(func(g Gomega) {
			g.Expect(k8sClient.Get(context.TODO(), key, pedals)).To(Succeed())
			pedals.Spec.Pressed = wanted
			g.Expect(k8sClient.Update(context.TODO(), pedals)).To(Succeed())
		}).Should(Succeed())

		By("watching the linkage move")
		expected := &playv1alpha1.Pedals{}
		Eventually(func(g Gomega) {
			g.Expect(k8sClient.Get(context.TODO(), key, expected)).To(Succeed())
			g.Expect(expected.Status.LinkagePosition).To(Equal(linkageWanted))
		}).Should(Succeed())

		By("checking that the rudder moved")
		rudderExpected := &playv1alpha1.Rudder{}
		Eventually(func(g Gomega) {
			g.Expect(k8sClient.Get(context.TODO(), key, rudderExpected)).To(Succeed())
			g.Expect(rudderExpected.Status.Position).To(Equal(linkageWanted))
		}).Should(Succeed())
	})
})
