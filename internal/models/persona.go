package models

type Persona struct {
	ID                   int    `json:"id"`
	Carrera              string `json:"carrera"`
	GeneroAccion         int    `json:"genero_accion"`
	GeneroCienciaFiccion int    `json:"genero_ciencia_ficcion"`
	GeneroComedia        int    `json:"genero_comedia"`
	GeneroTerror         int    `json:"genero_terror"`
	GeneroDocumental     int    `json:"genero_documental"`
	GeneroRomance        int    `json:"genero_romance"`
	GeneroMusicales      int    `json:"genero_musicales"`

	Poo                 float64 `json:"poo"`
	Ctd                 float64 `json:"ctd"`
	CalculoMultivariado float64 `json:"calculo_multivariado"`

	IngenieriaSoftware *float64 `json:"ingenieria_software"`
	BasesDatos         *float64 `json:"bases_datos"`
	ControlAnalogo     *float64 `json:"control_analogo"`
	CircuitosDigitales *float64 `json:"circuitos_digitales"`

	Promedio float64 `json:"promedio"`
}
