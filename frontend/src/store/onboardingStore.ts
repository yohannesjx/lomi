import { create } from 'zustand';
import { OnboardingService, OnboardingStatus } from '../api/onboarding';

interface OnboardingState {
    onboardingStep: number;
    onboardingCompleted: boolean;
    progress: number;
    isLoading: boolean;
    fetchStatus: () => Promise<void>;
    updateStep: (step: number, completed?: boolean) => Promise<void>;
    reset: () => void;
}

export const useOnboardingStore = create<OnboardingState>((set, get) => ({
    onboardingStep: 0,
    onboardingCompleted: false,
    progress: 0,
    isLoading: false,

    fetchStatus: async () => {
        try {
            set({ isLoading: true });
            const status = await OnboardingService.getStatus();
            set({
                onboardingStep: status.onboarding_step,
                onboardingCompleted: status.onboarding_completed,
                progress: status.progress,
                isLoading: false,
            });
            console.log('✅ Onboarding status fetched:', status);
        } catch (error) {
            console.error('❌ Failed to fetch onboarding status:', error);
            set({ isLoading: false });
        }
    },

    updateStep: async (step: number, completed?: boolean) => {
        try {
            console.log(`📤 Updating onboarding step to ${step}...`);
            const status = await OnboardingService.updateProgress(step, completed);
            set({
                onboardingStep: status.onboarding_step,
                onboardingCompleted: status.onboarding_completed,
                progress: status.progress,
            });
            console.log(`✅ Onboarding step updated: ${step}, completed: ${status.onboarding_completed}`);
            return status;
        } catch (error: any) {
            console.error('❌ Failed to update onboarding step:', error);
            console.error('Error details:', {
                message: error?.message,
                response: error?.response?.data,
                status: error?.response?.status,
            });
            throw error;
        }
    },

    reset: () => {
        set({
            onboardingStep: 0,
            onboardingCompleted: false,
            progress: 0,
        });
    },
}));

