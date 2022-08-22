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
	"reflect"
	"strings"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	playv1alpha1 "github.com/roehrich-hpe/airplane-sim/api/v1alpha1"
)

var _ = Describe("Airplane unit tests for initial population", func() {

	It("is invalid if the spec is not specified", func() {
		airplane := &playv1alpha1.Airplane{
			ObjectMeta: metav1.ObjectMeta{
				Name:      uuid.New().String()[0:8],
				Namespace: corev1.NamespaceDefault,
			},
		}

		Expect(k8sClient.Create(context.TODO(), airplane)).ToNot(Succeed())
	})

	It("is invalid if the spec.tailNumber is not specified", func() {
		airplane := &playv1alpha1.Airplane{
			ObjectMeta: metav1.ObjectMeta{
				Name:      uuid.New().String()[0:8],
				Namespace: corev1.NamespaceDefault,
			},
			Spec: playv1alpha1.AirplaneSpec{},
		}

		Expect(k8sClient.Create(context.TODO(), airplane)).ToNot(Succeed())
	})

	planeWithTailNumber := func(tailNumber string) *playv1alpha1.Airplane {
		return &playv1alpha1.Airplane{
			ObjectMeta: metav1.ObjectMeta{
				Name:      uuid.New().String()[0:8],
				Namespace: corev1.NamespaceDefault,
			},
			Spec: playv1alpha1.AirplaneSpec{
				TailNumber: tailNumber,
			},
		}
	}

	Describe("check tail number values", func() {
		var airplane *playv1alpha1.Airplane

		DescribeTable("verify invalid tail numbers",
			func(tailNumber string, isValid bool) {
				airplane = planeWithTailNumber(tailNumber)
				if isValid {
					Expect(k8sClient.Create(context.TODO(), airplane)).To(Succeed())
					Expect(k8sClient.Delete(context.TODO(), airplane)).To(Succeed())
				} else {
					Expect(k8sClient.Create(context.TODO(), airplane)).ToNot(Succeed())
				}
			},
			Entry("when empty", "", false),
			Entry("when lowercase N", "n123AB", false),
			Entry("when too short", "N1234", false),
			Entry("when no leading N", "z123AB", false),
			Entry("when some lowercase", "N123ab", false),
			Entry("when too long", "N123ABC", false),
			Entry("when good", "N123BC", true),
			Entry("when good", "N901NV", true),
		)
	})
})

var _ = Describe("Airplane unit tests", func() {

	var (
		key          types.NamespacedName
		ckey         types.NamespacedName
		ucTailNumber string
		tailNumber   string
		airplane     *playv1alpha1.Airplane
	)

	BeforeEach(func() {
		key = types.NamespacedName{
			Name:      uuid.New().String()[0:8],
			Namespace: corev1.NamespaceDefault,
		}

		ucTailNumber = "N" + strings.ToUpper(uuid.New().String()[0:5])
		tailNumber = strings.ToLower(ucTailNumber)

		airplane = &playv1alpha1.Airplane{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key.Name,
				Namespace: key.Namespace,
			},
			Spec: playv1alpha1.AirplaneSpec{
				TailNumber: ucTailNumber,
			},
		}

		Expect(k8sClient.Create(context.TODO(), airplane)).To(Succeed())

		// Keep a key that is used to find the components.
		ckey = types.NamespacedName{
			Name:      tailNumber,
			Namespace: airplane.GetNamespace(),
		}
	})

	AfterEach(func() {
		// The airplane should appear as the owner of its components.
		// In the test control plane the components won't be cleaned
		// up by garbage collection, so instead we settle for verifing
		// that they have the owner reference that would be used by
		// garbage collection.
		Expect(k8sClient.Get(context.TODO(), key, airplane)).To(Succeed())
		var truePtr bool = true
		expectedOwnerReference := v1.OwnerReference{
			Kind:               reflect.TypeOf(*airplane).Name(),
			APIVersion:         airplane.APIVersion,
			UID:                airplane.GetUID(),
			Name:               airplane.GetName(),
			Controller:         &truePtr,
			BlockOwnerDeletion: &truePtr,
		}

		pedals := &playv1alpha1.Pedals{}
		Expect(k8sClient.Get(context.TODO(), ckey, pedals)).To(Succeed())
		Expect(pedals.GetOwnerReferences()).To(HaveLen(1))
		Expect(pedals.GetOwnerReferences()).To(ContainElement(expectedOwnerReference))

		rudder := &playv1alpha1.Rudder{}
		Expect(k8sClient.Get(context.TODO(), ckey, rudder)).To(Succeed())
		Expect(rudder.GetOwnerReferences()).To(HaveLen(1))
		Expect(rudder.GetOwnerReferences()).To(ContainElement(expectedOwnerReference))

		// Now delete the airplane.
		Expect(k8sClient.Delete(context.TODO(), airplane)).To(Succeed())

		// Test env limitation: Wait until the cached object is gone.
		expected := &playv1alpha1.Airplane{}
		Eventually(func() error {
			return k8sClient.Get(context.TODO(), key, expected)
		}).ShouldNot(Succeed())
	})

	It("Creates pedals and rudder", func() {
		Eventually(func(g Gomega) {
			g.Expect(k8sClient.Get(context.TODO(), key, airplane)).To(Succeed())
			g.Expect(airplane.Status.Pedals.Name).To(Equal(tailNumber))
			g.Expect(airplane.Status.Rudder.Name).To(Equal(tailNumber))
		}).Should(Succeed())

		pedals := &playv1alpha1.Pedals{}
		Expect(k8sClient.Get(context.TODO(), ckey, pedals)).To(Succeed())
		rudder := &playv1alpha1.Rudder{}
		Expect(k8sClient.Get(context.TODO(), ckey, rudder)).To(Succeed())
	})
})
