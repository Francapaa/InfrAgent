"use client";

import { useAuth } from "@/app/hooks/useAuth";
import { LoadingState } from "./LoadingState";
import { OnboardingForm } from "./OnboardingForm";
import { SuccessView } from "./SuccessView";
import { useOnboarding } from "../hooks/useOnboarding";

export function OnboardingContainer() {
  const { isLoading, user } = useAuth();
  const { state, setField, submitForm, goToDashboard } = useOnboarding();

  if (isLoading) {
    return <LoadingState message="Verificando autenticación..." />;
  }

  if (state.success && state.profileData) {
    return <SuccessView profileData={state.profileData} onGoToDashboard={goToDashboard} />;
  }

  return (
    <OnboardingForm
      companyName={state.formData.companyName}
      webhookUrl={state.formData.webhookUrl}
      isSubmitting={state.isSubmitting}
      error={state.error}
      onCompanyNameChange={(value) => setField("companyName", value)}
      onWebhookUrlChange={(value) => setField("webhookUrl", value)}
      onSubmit={submitForm}
    />
  );
}
