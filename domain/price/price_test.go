package price_test

import (
	"encoding/json"
	"testing"

	domainerrors "github.com/mytheresa/go-hiring-challenge/domain/errors"
	"github.com/mytheresa/go-hiring-challenge/domain/price"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrice(t *testing.T) {
	t.Run("New wraps a decimal value", func(t *testing.T) {
		d := decimal.NewFromFloat(19.99)
		p := price.New(d)
		assert.True(t, d.Equal(p.Decimal()))
	})

	t.Run("Parse", func(t *testing.T) {
		cases := []struct {
			name    string
			input   string
			wantErr bool
		}{
			{name: "valid decimal string", input: "19.99"},
			{name: "invalid decimal string", input: "invalid", wantErr: true},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				p, err := price.Parse(c.input)
				if c.wantErr {
					assert.Error(t, err)
					assert.ErrorIs(t, err, domainerrors.ErrorDomainValidation)
					return
				}
				require.NoError(t, err)
				assert.Equal(t, c.input, p.String())
			})
		}
	})

	t.Run("Equal", func(t *testing.T) {
		cases := []struct {
			name  string
			a     string
			b     string
			equal bool
		}{
			{name: "equal values", a: "19.99", b: "19.99", equal: true},
			{name: "equal values different representation", a: "19.90", b: "19.9", equal: true},
			{name: "different values", a: "19.99", b: "20.00", equal: false},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				a, err := price.Parse(c.a)
				require.NoError(t, err)
				b, err := price.Parse(c.b)
				require.NoError(t, err)

				assert.Equal(t, c.equal, a.Equal(b))
			})
		}
	})

	t.Run("String formats the decimal value", func(t *testing.T) {
		p, err := price.Parse("19.99")
		require.NoError(t, err)
		assert.Equal(t, "19.99", p.String())
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		p, err := price.Parse("19.99")
		require.NoError(t, err)

		data, err := json.Marshal(p)
		require.NoError(t, err)
		assert.Equal(t, `"19.99"`, string(data))
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		cases := []struct {
			name    string
			input   string
			wantErr bool
			want    string
		}{
			{name: "valid JSON number", input: "19.99", want: "19.99"},
			{name: "invalid JSON", input: "not-a-number", wantErr: true},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				var p price.Price
				err := json.Unmarshal([]byte(c.input), &p)
				if c.wantErr {
					assert.Error(t, err)
					return
				}
				require.NoError(t, err)
				assert.Equal(t, c.want, p.String())
			})
		}
	})
}
