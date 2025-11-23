import React from 'react';
import { View, Text, StyleSheet, Animated } from 'react-native';
import { COLORS, SPACING } from '../../theme/colors';

interface ProgressBarProps {
    currentStep: number;
    totalSteps?: number;
}

export const ProgressBar: React.FC<ProgressBarProps> = ({
    currentStep,
    totalSteps = 8
}) => {
    const progress = (currentStep / totalSteps) * 100;
    const progressAnim = React.useRef(new Animated.Value(0)).current;

    React.useEffect(() => {
        Animated.timing(progressAnim, {
            toValue: progress,
            duration: 300,
            useNativeDriver: false,
        }).start();
    }, [progress]);

    return (
        <View style={styles.container}>
            <View style={styles.progressBarBg}>
                <Animated.View
                    style={[
                        styles.progressBarFill,
                        {
                            width: progressAnim.interpolate({
                                inputRange: [0, 100],
                                outputRange: ['0%', '100%'],
                            }),
                        },
                    ]}
                />
            </View>
            <Text style={styles.stepText}>
                Step {currentStep} of {totalSteps}
            </Text>
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        paddingHorizontal: SPACING.l,
        paddingVertical: SPACING.m,
    },
    progressBarBg: {
        height: 4,
        backgroundColor: COLORS.surfaceHighlight,
        borderRadius: 2,
        overflow: 'hidden',
        marginBottom: SPACING.xs,
    },
    progressBarFill: {
        height: '100%',
        backgroundColor: COLORS.primary,
        borderRadius: 2,
    },
    stepText: {
        fontSize: 12,
        color: COLORS.textSecondary,
        textAlign: 'center',
    },
});
