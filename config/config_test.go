package config

import (
	"testing"

	. "github.com/onsi/gomega"
)

type SubConfig struct {
	FiledA *string
}

func (c *SubConfig) Normalize() error {
	if c.FiledA == nil {
		s := "default-in-subconfig"
		c.FiledA = &s
	} else {
		s := *c.FiledA + "-modified"
		c.FiledA = &s
	}
	return nil
}

type Config struct {
	FieldA string
	FieldB *int

	SubConfig *SubConfig
}

func (c *Config) Normalize() error {
	var err error

	if c.FieldA == "" {
		c.FieldA = "default-filedA-in-config"
	} else {
		c.FieldA += "-modified"
	}

	if c.FieldB == nil {
		i := 1
		c.FieldB = &i
	}

	c.SubConfig, err = NormalizeOrDefault(c.SubConfig)
	if err != nil {
		return err
	}

	return nil
}

func TestConfig(t *testing.T) {
	cases := map[string]struct {
		input     *Config
		expectErr error
		expect    *Config
	}{
		"return default value if input is nil": {
			input: nil,
			expect: &Config{
				FieldA: "default-filedA-in-config",
				FieldB: func() *int {
					i := 1
					return &i
				}(),
				SubConfig: &SubConfig{
					FiledA: func() *string {
						s := "default-in-subconfig"
						return &s
					}(),
				},
			},
		},
		"complete default value if input is not nil": {
			input: &Config{
				FieldA: "set-fieldA",
				FieldB: func() *int {
					i := 2
					return &i
				}(),
				SubConfig: &SubConfig{
					FiledA: func() *string {
						s := "set-in-subconfig"
						return &s
					}(),
				},
			},
			expect: &Config{
				FieldA: "set-fieldA-modified",
				FieldB: func() *int {
					i := 2
					return &i
				}(),
				SubConfig: &SubConfig{
					FiledA: func() *string {
						s := "set-in-subconfig-modified"
						return &s
					}(),
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			actual, err := NormalizeOrDefault(tc.input)
			if tc.expectErr != nil {
				g.Expect(err).To(Equal(tc.expectErr))
			} else {
				g.Expect(err).To(BeNil())
			}
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}
