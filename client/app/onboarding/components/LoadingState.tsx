"use client";

interface LoadingStateProps {
  message?: string;
}

export function LoadingState({ message = "Cargando..." }: LoadingStateProps) {
  return (
    <div className="min-h-screen bg-black flex items-center justify-center">
      <div className="text-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-yellow-400 mx-auto mb-4"></div>
        <p className="text-gray-400">{message}</p>
      </div>
    </div>
  );
}
