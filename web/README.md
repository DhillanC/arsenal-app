# web

## ¿Qué es?

Frontend de la aplicación. Server-rendered con Go templates + HTMX + Tailwind CSS. Tema visual basado en DCS Web.

## Subcarpetas

| Carpeta | Contenido |
|---------|-----------|
| `templates/` | Templates HTML (Go html/template) |
| `static/css/` | Estilos CSS con paleta DCS |
| `static/js/` | JavaScript mínimo (theme toggle, HTMX boosts) |

## Templates principales

- `base.html` — Layout base con dark/light mode
- `layout.html` — Estructura común de página
- `index.html` — Dashboard principal
- `replica_list.html` — Lista de réplicas
- `replica_detail.html` — Ficha de réplica con tabs
- `replica_form.html` — Formulario crear/editar
- `document_list.html` — Lista de documentos
- `mantenimiento_list.html` — Próximos mantenimientos
- `error.html` — Página de error

## Tema DCS

| Token | Valor | Uso |
|-------|-------|-----|
| Gold Primary | `#b88834` | Botones, acentos |
| Gold Light | `#fdf3aa` | Textos destacados |
| Near-black | `#131110` | Fondo dark mode |
| Cream | `#f9f6f0` | Fondo light mode |
| Teal | `#5DC8D2` | Acento dark mode |

## Dark/Light Mode

- Toggle guarda preferencia en `localStorage` (key: `theme`)
- Clase `.dark` en `:root` activa modo oscuro
- Transiciones suaves entre modos

## PWA

- `manifest.json` con iconos y colores
- Service worker placeholder
- Meta tags para mobile
