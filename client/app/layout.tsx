import '../app/globals.css'
import { inter, jetbrainsMono } from '../app/fonts/fonts'
import type { Metadata, Viewport } from 'next'

export const metadata: Metadata = {
  title: 'InfrAgent - Autonomous AI Infrastructure Agent',
  description:
    'InfrAgent monitors your infrastructure 24/7, detects issues in real-time, and autonomously resolves them with AI-powered decision making.',
}

export const viewport: Viewport = {
  themeColor: '#000000',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="es">
      <body
        className={`${inter.variable} ${jetbrainsMono.variable} font-sans antialiased`}
      >
        {children}
      </body>
    </html>
  )
}
