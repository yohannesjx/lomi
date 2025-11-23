import { Platform } from 'react-native';
import * as Haptics from 'expo-haptics';

/**
 * Cross-platform haptics utility
 * Safe to call on web (no-op)
 */
export const triggerHaptic = {
    impact: (style: Haptics.ImpactFeedbackStyle = Haptics.ImpactFeedbackStyle.Light) => {
        if (Platform.OS !== 'web') {
            Haptics.impactAsync(style).catch(() => {
                // Silently fail if haptics not available
            });
        }
    },
    notification: (type: Haptics.NotificationFeedbackType = Haptics.NotificationFeedbackType.Success) => {
        if (Platform.OS !== 'web') {
            Haptics.notificationAsync(type).catch(() => {
                // Silently fail if haptics not available
            });
        }
    },
};

