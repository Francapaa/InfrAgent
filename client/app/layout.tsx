import '../app/globals.css' // Asegúrate de importar los estilos de Tailwind
import {inter} from '../app/fonts/fonts'
export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="es">
      <body className={`${inter.className} antialised`}>{children}</body>
    </html>
  )
}