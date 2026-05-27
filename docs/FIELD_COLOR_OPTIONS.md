# Opciones de Color para Campos de Formulario

## Contexto

Los campos de texto, selects, textareas y file inputs necesitan colores explícitos por modo visual.

Problemas actuales:

- En modo claro, los campos blancos se pierden contra el fondo.
- En modo oscuro, el fondo del campo puede verse, pero el texto queda demasiado claro o con bajo contraste durante la escritura.
- La solución debe definir fondo, borde, texto, placeholder y estado focus para ambos modos.

No aplicar cambios de CSS todavía. Este documento solo registra opciones para decidir.

## Opción 1: Neutro profesional

Recomendada por legibilidad y bajo riesgo visual.

### Modo claro

- Fondo: `#f3f4f6`
- Borde: `#d1d5db`
- Texto: `#111827`
- Placeholder: `#6b7280`

### Modo oscuro

- Fondo: `#1f2937`
- Borde: `#4b5563`
- Texto: `#f9fafb`
- Placeholder: `#9ca3af`

### Focus

- Borde: `#b88834`
- Halo: `rgba(184,136,52,0.18)`

### Comentario

Es la opción más limpia para una app operativa. No compite con el dorado DCS ni cambia demasiado la identidad visual.

## Opción 2: Azul claro y slate

Más moderna, con mayor separación visual en modo claro.

### Modo claro

- Fondo: `#eff6ff`
- Borde: `#bfdbfe`
- Texto: `#0f172a`
- Placeholder: `#64748b`

### Modo oscuro

- Fondo: `#172033`
- Borde: `#334155`
- Texto: `#f8fafc`
- Placeholder: `#94a3b8`

### Focus

- Borde: `#5DC8D2`
- Halo: `rgba(93,200,210,0.18)`

### Comentario

Da buena legibilidad y un look más tecnológico, pero puede competir visualmente con el dorado del tema DCS.

## Opción 3: Crema y carbón, alineado con DCS

Más integrada al tema actual.

### Modo claro

- Fondo: `#f5efe3`
- Borde: `#d8c7a3`
- Texto: `#1f1a14`
- Placeholder: `#7c6f60`

### Modo oscuro

- Fondo: `#241f1a`
- Borde: `#5b4630`
- Texto: `#f7efe4`
- Placeholder: `#b8a895`

### Focus

- Borde: `#b88834`
- Halo: `rgba(184,136,52,0.20)`

### Comentario

Se siente más propio del tema DCS, pero hay que evitar que la interfaz termine demasiado beige o monocromática.

## Opción 4: Blanco controlado y gris oscuro

Menor cambio visual sobre el diseño actual.

### Modo claro

- Fondo: `#ffffff`
- Borde: `#cbd5e1`
- Texto: `#0f172a`
- Placeholder: `#64748b`

### Modo oscuro

- Fondo: `#111827`
- Borde: `#374151`
- Texto: `#f9fafb`
- Placeholder: `#9ca3af`

### Focus

- Borde: `#b88834`
- Halo: `rgba(184,136,52,0.18)`

### Comentario

Conserva inputs blancos en claro, pero mejora contraste mediante borde y texto. Es menos efectivo si el problema principal es que el campo blanco se pierde en la página.

## Recomendación inicial

Usar la opción 1: Neutro profesional.

Razones:

- Tiene contraste fuerte en ambos modos.
- No compite con la paleta DCS.
- Funciona bien para inputs, selects, textareas y file inputs.
- El focus dorado mantiene la identidad visual sin saturar la interfaz.

## Pendientes antes de implementar

- Revisar todos los campos existentes: `input`, `select`, `textarea`, `input[type="file"]`.
- Definir estilos de `disabled`, `readonly`, `error` y `success`.
- Verificar contraste en modo claro y oscuro.
- Probar en mobile, especialmente inputs dentro de cards.
