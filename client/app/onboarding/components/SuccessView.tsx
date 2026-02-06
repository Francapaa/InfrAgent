"use client";

import { CompleteProfileResponse } from "@/lib/api/types";

interface SuccessViewProps {
  profileData: CompleteProfileResponse;
  onGoToDashboard: () => void;
}

function copyToClipboard(text: string): void {
  navigator.clipboard.writeText(text);
}

interface CredentialFieldProps {
  label: string;
  value: string;
}

function CredentialField({ label, value }: CredentialFieldProps) {
  return (
    <div>
      <label className="block text-sm font-medium text-gray-300 mb-2">
        {label}
      </label>
      <div className="flex gap-2">
        <input
          type="text"
          value={value}
          readOnly
          className="flex-1 px-4 py-3 bg-black border border-zinc-700 rounded-lg text-white font-mono text-sm"
        />
        <button
          onClick={() => copyToClipboard(value)}
          className="px-4 py-3 bg-yellow-400 text-black font-semibold rounded-lg hover:bg-yellow-300 transition"
        >
          Copiar
        </button>
      </div>
    </div>
  );
}

export function SuccessView({ profileData, onGoToDashboard }: SuccessViewProps) {
  return (
    <div className="min-h-screen bg-black flex items-center justify-center p-4">
      <div className="w-full max-w-2xl">
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-yellow-400 mb-2">InfrAgent</h1>
          <p className="text-gray-400">¡Perfil completado exitosamente!</p>
        </div>

        <div className="bg-zinc-900 rounded-2xl p-8 shadow-2xl border border-zinc-800">
          <div className="mb-6">
            <div className="bg-green-500/10 border border-green-500/20 rounded-lg p-4 mb-6">
              <p className="text-green-400 text-center">
                <strong>¡Importante!</strong> Guarda estas credenciales en un lugar seguro. No
                podrás verlas nuevamente.
              </p>
            </div>
          </div>

          <div className="space-y-6">
            <CredentialField label="API Key" value={profileData.api_key} />

            <CredentialField label="Webhook Secret" value={profileData.webhook_secret} />

            <div className="bg-zinc-800/50 rounded-lg p-4 mt-6">
              <h3 className="text-yellow-400 font-semibold mb-2">Siguientes pasos:</h3>
              <ul className="text-gray-400 text-sm space-y-2 list-disc list-inside">
                <li>Guarda tu API Key para autenticar tus solicitudes</li>
                <li>Configura el Webhook Secret en tu servidor para verificar las firmas</li>
                <li>Revisa la documentación para integrar el SDK</li>
              </ul>
            </div>

            <button
              onClick={onGoToDashboard}
              className="w-full bg-yellow-400 text-black font-semibold py-3 rounded-lg hover:bg-yellow-300 transition duration-200 shadow-lg hover:shadow-yellow-400/50 mt-6"
            >
              Ir al Dashboard
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
