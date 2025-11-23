import React, { useEffect, useRef } from 'react';
import { View, Text, StyleSheet, Animated, Dimensions } from 'react-native';
import { COLORS, SPACING, SIZES } from '../../theme/colors';

const { width, height } = Dimensions.get('window');

interface GiftAnimationProps {
    visible: boolean;
    gift: {
        name_en?: string;
        name_am?: string;
        icon_url?: string;
        animation_url?: string;
    };
    onComplete: () => void;
}

export const GiftAnimation: React.FC<GiftAnimationProps> = ({ visible, gift, onComplete }) => {
    const scaleAnim = useRef(new Animated.Value(0)).current;
    const rotateAnim = useRef(new Animated.Value(0)).current;
    const fadeAnim = useRef(new Animated.Value(0)).current;
    const sparkleAnim = useRef(new Animated.Value(0)).current;

    useEffect(() => {
        if (visible) {
            // Reset animations
            scaleAnim.setValue(0);
            rotateAnim.setValue(0);
            fadeAnim.setValue(0);
            sparkleAnim.setValue(0);

            // Create sparkle animation
            const sparkleAnimation = Animated.loop(
                Animated.sequence([
                    Animated.timing(sparkleAnim, {
                        toValue: 1,
                        duration: 500,
                        useNativeDriver: true,
                    }),
                    Animated.timing(sparkleAnim, {
                        toValue: 0,
                        duration: 500,
                        useNativeDriver: true,
                    }),
                ])
            );

            // Main animation sequence
            Animated.parallel([
                // Scale up with bounce
                Animated.sequence([
                    Animated.spring(scaleAnim, {
                        toValue: 1.2,
                        friction: 3,
                        tension: 40,
                        useNativeDriver: true,
                    }),
                    Animated.spring(scaleAnim, {
                        toValue: 1,
                        friction: 3,
                        tension: 40,
                        useNativeDriver: true,
                    }),
                ]),
                // Rotate
                Animated.timing(rotateAnim, {
                    toValue: 1,
                    duration: 1000,
                    useNativeDriver: true,
                }),
                // Fade in
                Animated.timing(fadeAnim, {
                    toValue: 1,
                    duration: 300,
                    useNativeDriver: true,
                }),
                // Sparkle
                sparkleAnimation,
            ]).start();

            // Auto-close after 3 seconds
            const timer = setTimeout(() => {
                Animated.parallel([
                    Animated.timing(scaleAnim, {
                        toValue: 0,
                        duration: 300,
                        useNativeDriver: true,
                    }),
                    Animated.timing(fadeAnim, {
                        toValue: 0,
                        duration: 300,
                        useNativeDriver: true,
                    }),
                ]).start(() => {
                    onComplete();
                });
            }, 3000);

            return () => clearTimeout(timer);
        }
    }, [visible]);

    if (!visible) return null;

    const rotate = rotateAnim.interpolate({
        inputRange: [0, 1],
        outputRange: ['0deg', '360deg'],
    });

    const sparkleOpacity = sparkleAnim.interpolate({
        inputRange: [0, 0.5, 1],
        outputRange: [0.3, 1, 0.3],
    });

    return (
        <View style={styles.overlay} pointerEvents="none">
            <Animated.View
                style={[
                    styles.container,
                    {
                        opacity: fadeAnim,
                        transform: [{ scale: scaleAnim }, { rotate }],
                    },
                ]}
            >
                {/* Sparkle effects */}
                {[...Array(8)].map((_, i) => {
                    const angle = (i * 360) / 8;
                    const radius = 150;
                    const x = Math.cos((angle * Math.PI) / 180) * radius;
                    const y = Math.sin((angle * Math.PI) / 180) * radius;

                    return (
                        <Animated.View
                            key={i}
                            style={[
                                styles.sparkle,
                                {
                                    left: width / 2 + x - 10,
                                    top: height / 2 + y - 10,
                                    opacity: sparkleOpacity,
                                },
                            ]}
                        >
                            <Text style={styles.sparkleText}>✨</Text>
                        </Animated.View>
                    );
                })}

                {/* Gift Icon */}
                <View style={styles.giftContainer}>
                    {gift.icon_url ? (
                        <Text style={styles.giftEmoji}>🎁</Text>
                    ) : (
                        <Text style={styles.giftEmoji}>🎁</Text>
                    )}
                </View>

                {/* Gift Name */}
                <Text style={styles.giftName}>{gift.name_en || gift.name_am || 'Gift'}</Text>
                <Text style={styles.giftMessage}>You received a gift! 🎉</Text>
            </Animated.View>
        </View>
    );
};

const styles = StyleSheet.create({
    overlay: {
        ...StyleSheet.absoluteFillObject,
        backgroundColor: 'rgba(0, 0, 0, 0.8)',
        justifyContent: 'center',
        alignItems: 'center',
        zIndex: 9999,
    },
    container: {
        alignItems: 'center',
        justifyContent: 'center',
    },
    sparkle: {
        position: 'absolute',
        width: 20,
        height: 20,
    },
    sparkleText: {
        fontSize: 20,
    },
    giftContainer: {
        width: 200,
        height: 200,
        borderRadius: 100,
        backgroundColor: 'rgba(255, 215, 0, 0.2)',
        alignItems: 'center',
        justifyContent: 'center',
        marginBottom: SPACING.l,
        borderWidth: 3,
        borderColor: COLORS.gold,
    },
    giftEmoji: {
        fontSize: 100,
    },
    giftName: {
        fontSize: 32,
        fontWeight: 'bold',
        color: COLORS.gold,
        marginBottom: SPACING.s,
        textAlign: 'center',
    },
    giftMessage: {
        fontSize: 20,
        color: COLORS.background,
        textAlign: 'center',
    },
});

