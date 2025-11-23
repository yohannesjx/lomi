import React, { useEffect } from 'react';
import { View, Text, StyleSheet, Modal, TouchableOpacity, Image, Dimensions } from 'react-native';
import Animated, {
    useSharedValue,
    useAnimatedStyle,
    withSpring,
    withSequence,
    withTiming,
    withRepeat,
    interpolate,
} from 'react-native-reanimated';
import { LinearGradient } from 'expo-linear-gradient';
import { COLORS, SPACING, SIZES } from '../../theme/colors';

const { width, height } = Dimensions.get('window');

interface MatchModalProps {
    visible: boolean;
    matchedUser: {
        id: string;
        name: string;
        photos?: Array<{ url: string }>;
    };
    onClose: () => void;
    onSendMessage: () => void;
    onKeepSwiping: () => void;
}

export const MatchModal: React.FC<MatchModalProps> = ({
    visible,
    matchedUser,
    onClose,
    onSendMessage,
    onKeepSwiping,
}) => {
    const scale = useSharedValue(0);
    const opacity = useSharedValue(0);
    const avatar1Scale = useSharedValue(0);
    const avatar2Scale = useSharedValue(0);
    const confettiOpacity = useSharedValue(1);

    useEffect(() => {
        if (visible) {
            scale.value = withSpring(1, { damping: 12, stiffness: 100 });
            opacity.value = withTiming(1, { duration: 300 });
            avatar1Scale.value = withSequence(
                withTiming(0, { duration: 0 }),
                withSpring(1, { damping: 8, stiffness: 120 }),
                withSpring(1.1, { damping: 6 }),
                withSpring(1, { damping: 8 })
            );
            avatar2Scale.value = withSequence(
                withTiming(0, { duration: 0 }),
                withTiming(0, { duration: 200 }),
                withSpring(1, { damping: 8, stiffness: 120 }),
                withSpring(1.1, { damping: 6 }),
                withSpring(1, { damping: 8 })
            );
            confettiOpacity.value = withSequence(
                withTiming(1, { duration: 500 }),
                withTiming(1, { duration: 2000 }),
                withTiming(0, { duration: 500 })
            );
        } else {
            scale.value = 0;
            opacity.value = 0;
            avatar1Scale.value = 0;
            avatar2Scale.value = 0;
        }
    }, [visible]);

    const containerStyle = useAnimatedStyle(() => ({
        transform: [{ scale: scale.value }],
        opacity: opacity.value,
    }));

    const avatar1Style = useAnimatedStyle(() => ({
        transform: [{ scale: avatar1Scale.value }],
    }));

    const avatar2Style = useAnimatedStyle(() => ({
        transform: [{ scale: avatar2Scale.value }],
    }));

    const confettiStyle = useAnimatedStyle(() => ({
        opacity: confettiOpacity.value,
    }));

    const avatarUrl = matchedUser.photos?.[0]?.url || '';

    return (
        <Modal
            visible={visible}
            transparent
            animationType="none"
            onRequestClose={onClose}
        >
            <View style={styles.overlay}>
                {/* Confetti effect */}
                <Animated.View style={[styles.confettiContainer, confettiStyle]} pointerEvents="none">
                    <Text style={styles.confetti}>🎉</Text>
                    <Text style={[styles.confetti, { top: '20%', left: '10%' }]}>✨</Text>
                    <Text style={[styles.confetti, { top: '30%', right: '15%' }]}>💚</Text>
                    <Text style={[styles.confetti, { top: '50%', left: '20%' }]}>🎊</Text>
                    <Text style={[styles.confetti, { top: '60%', right: '10%' }]}>⭐</Text>
                    <Text style={[styles.confetti, { top: '70%', left: '15%' }]}>💫</Text>
                </Animated.View>

                <Animated.View style={[styles.container, containerStyle]}>
                    <LinearGradient
                        colors={['rgba(0,0,0,0.95)', 'rgba(0,0,0,0.98)']}
                        style={styles.gradient}
                    >
                        <Text style={styles.title}>It's a Match!</Text>
                        <Text style={styles.subtitle}>You and {matchedUser.name} liked each other</Text>

                        <View style={styles.avatarsContainer}>
                            <Animated.View style={avatar1Style}>
                                <View style={styles.avatarPlaceholder}>
                                    <Text style={styles.avatarText}>You</Text>
                                </View>
                            </Animated.View>

                            <View style={styles.heartContainer}>
                                <Text style={styles.heart}>💚</Text>
                            </View>

                            <Animated.View style={avatar2Style}>
                                {avatarUrl ? (
                                    <Image source={{ uri: avatarUrl }} style={styles.avatar} />
                                ) : (
                                    <View style={styles.avatarPlaceholder}>
                                        <Text style={styles.avatarText}>{matchedUser.name[0]}</Text>
                                    </View>
                                )}
                            </Animated.View>
                        </View>

                        <View style={styles.buttonsContainer}>
                            <TouchableOpacity
                                style={[styles.button, styles.primaryButton]}
                                onPress={onSendMessage}
                                activeOpacity={0.8}
                            >
                                <LinearGradient
                                    colors={[COLORS.primary, COLORS.primaryDark]}
                                    style={styles.buttonGradient}
                                >
                                    <Text style={styles.primaryButtonText}>Send a Message</Text>
                                </LinearGradient>
                            </TouchableOpacity>

                            <TouchableOpacity
                                style={[styles.button, styles.secondaryButton]}
                                onPress={onKeepSwiping}
                                activeOpacity={0.8}
                            >
                                <Text style={styles.secondaryButtonText}>Keep Swiping</Text>
                            </TouchableOpacity>
                        </View>
                    </LinearGradient>
                </Animated.View>
            </View>
        </Modal>
    );
};

const styles = StyleSheet.create({
    overlay: {
        flex: 1,
        backgroundColor: 'rgba(0,0,0,0.9)',
        justifyContent: 'center',
        alignItems: 'center',
    },
    confettiContainer: {
        position: 'absolute',
        width: '100%',
        height: '100%',
        zIndex: 1,
    },
    confetti: {
        position: 'absolute',
        fontSize: 40,
        top: '10%',
        right: '20%',
    },
    container: {
        width: width * 0.85,
        borderRadius: 24,
        overflow: 'hidden',
        zIndex: 2,
    },
    gradient: {
        padding: SPACING.xl,
        alignItems: 'center',
    },
    title: {
        fontSize: 42,
        fontWeight: '900',
        color: COLORS.primary,
        marginBottom: SPACING.s,
        textAlign: 'center',
    },
    subtitle: {
        fontSize: 18,
        color: COLORS.textSecondary,
        marginBottom: SPACING.xl,
        textAlign: 'center',
    },
    avatarsContainer: {
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'center',
        marginBottom: SPACING.xl,
        gap: SPACING.l,
    },
    avatar: {
        width: 100,
        height: 100,
        borderRadius: 50,
        borderWidth: 4,
        borderColor: COLORS.primary,
    },
    avatarPlaceholder: {
        width: 100,
        height: 100,
        borderRadius: 50,
        backgroundColor: COLORS.surface,
        borderWidth: 4,
        borderColor: COLORS.primary,
        alignItems: 'center',
        justifyContent: 'center',
    },
    avatarText: {
        fontSize: 36,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
    },
    heartContainer: {
        alignItems: 'center',
        justifyContent: 'center',
    },
    heart: {
        fontSize: 48,
    },
    buttonsContainer: {
        width: '100%',
        gap: SPACING.m,
    },
    button: {
        width: '100%',
        height: 56,
        borderRadius: 28,
        overflow: 'hidden',
    },
    buttonGradient: {
        flex: 1,
        alignItems: 'center',
        justifyContent: 'center',
    },
    primaryButton: {
        shadowColor: COLORS.primary,
        shadowOffset: { width: 0, height: 4 },
        shadowOpacity: 0.4,
        shadowRadius: 8,
        elevation: 8,
    },
    primaryButtonText: {
        fontSize: 18,
        fontWeight: 'bold',
        color: COLORS.background,
    },
    secondaryButton: {
        backgroundColor: COLORS.surface,
        borderWidth: 2,
        borderColor: COLORS.textSecondary,
        alignItems: 'center',
        justifyContent: 'center',
    },
    secondaryButtonText: {
        fontSize: 18,
        fontWeight: '600',
        color: COLORS.textPrimary,
    },
});

