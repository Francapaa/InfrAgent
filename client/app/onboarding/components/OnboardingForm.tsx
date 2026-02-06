"use client";

import { FormEvent } from "react";

interface OnboardingFormProps {
  companyName: string;
  webhookUrl: string;
  isSubmitting: boolean;
  error: string | null;
  onCompanyNameChange: (value: string) => void;
  onWebhookUrlChange: (value: string) => void;
  onSubmit: () => void;
}

export function OnboardingForm({
  companyName,
  webhookUrl,
  isSubmitting,
  error,
  onCompanyNameChange,
  onWebhookUrlChange,
  onSubmit,
}: OnboardingFormProps) {
  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    onSubmit();
  };

  return (
    <div className="min-h-screen bg-black flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-yellow-400 mb-2">InfrAgent</h1>
          <p className="text-gray-400 text-sm">Completa tu perfil para continuar</p>
        </div>

        <div className="bg-zinc-900 rounded-2xl p-8 shadow-2xl border border-zinc-800">
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label
                htmlFor="companyName"
                className="block text-sm font-medium text-gray-300 mb-2"
              >
                Nombre de la empresa o proyecto
              </label>
              <input
                id="companyName"
                type="text"
                value={companyName}
                onChange={(e) => onCompanyNameChange(e.target.value)}
                className="w-full px-4 py-3 bg-black border border-zinc-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-yellow-400 focus:border-transparent transition"
                placeholder="Mi Empresa"
                required
                disabled={isSubmitting}
              />
            </div>

            <div>
              <label
                htmlFor="webhookUrl"
                className="block text-sm font-medium text-gray-300 mb-2"
              >
                Webhook URL
              </label>
              <input
                id="webhookUrl"
                type="url"
                value={webhookUrl}
                onChange={(e) => onWebhookUrlChange(e.target.value)}
                className="w-full px-4 py-3 bg-black border border-zinc-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-yellow-400 focus:border-transparent transition"
                placeholder="https://tu-servidor.com/webhook"
                required
                disabled={isSubmitting}
              />
              <p className="text-gray-500 text-xs mt-1">Debe usar HTTPS</p>
            </div>

            {error && (
              <div className="bg-red-500/10 border border-red-500/20 rounded-lg p-3">
                <p className="text-red-400 text-sm text-center">{error}</p>
              </div>
            )}

            <button
              type="submit"
              disabled={isSubmitting}
              className="w-full bg-yellow-400 text-black font-semibold py-3 rounded-lg hover:bg-yellow-300 transition duration-200 shadow-lg hover:shadow-yellow-400/50 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isSubmitting ? "Guardando..." : "Completar perfil"}
            </button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-gray-500 text-xs">
              Este paso es obligatorio para usar la plataforma
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
