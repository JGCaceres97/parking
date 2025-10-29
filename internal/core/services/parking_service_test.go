package services

import (
	"testing"
	"time"
)

func TestCalculateCharge(t *testing.T) {
	tests := []struct {
		name           string
		duration       time.Duration
		rate           float64
		expectedHours  int
		expectedCharge float64
	}{
		// Caso 1: Menos de 1 minuto (se cobra 1 hora)
		{"Menos de un minuto", time.Minute * 0, 15.00, 1, 15.00},

		// Caso 2: 1 hora justa
		{"1 hora justa", time.Hour * 1, 15.00, 1, 15.00},

		// Caso 3: 1 hora y 29 minutos (Se queda en 1 hora)
		{"1h 29min (Normal)", time.Hour*1 + time.Minute*29, 15.00, 1, 15.00},

		// Caso 4: 1 hora y 30 minutos (Redondea a 2 horas)
		{"1h 30min (Especial)", time.Hour*1 + time.Minute*30, 5.00, 2, 10.00},

		// Caso 5: 2 horas y 1 minuto (Se queda en 2 horas)
		{"2h 1min (Normal)", time.Hour*2 + time.Minute*1, 15.00, 2, 30.00},

		// Caso 6: Motocicleta exenta (3 horas y 45 minutos)
		{"Motocicleta Exenta", time.Hour*3 + time.Minute*45, 0.00, 4, 0.00},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entryTime := time.Now()
			exitTime := entryTime.Add(tt.duration)

			hours, charge := calculateCharge(entryTime, exitTime, tt.rate)

			if hours != tt.expectedHours {
				t.Errorf("Horas cobradas incorrectas. Esperado: %d, Obtenido: %d", tt.expectedHours, hours)
			}

			if charge != tt.expectedCharge {
				t.Errorf("Cobro total incorrecto. Esperado: %.2f, Obtenido: %.2f", tt.expectedCharge, charge)
			}
		})
	}
}
