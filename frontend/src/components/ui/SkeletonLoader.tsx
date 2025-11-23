import React, { useEffect } from 'react';
import { View, StyleSheet, Animated } from 'react-native';
import { COLORS, SPACING, SIZES } from '../../theme/colors';

interface SkeletonProps {
    width?: number | string;
    height?: number;
    borderRadius?: number;
    style?: any;
}

export const SkeletonLoader: React.FC<SkeletonProps> = ({
    width = '100%',
    height = 20,
    borderRadius = SIZES.radiusS,
    style,
}) => {
    const animatedValue = new Animated.Value(0);

    useEffect(() => {
        Animated.loop(
            Animated.sequence([
                Animated.timing(animatedValue, {
                    toValue: 1,
                    duration: 1000,
                    useNativeDriver: true,
                }),
                Animated.timing(animatedValue, {
                    toValue: 0,
                    duration: 1000,
                    useNativeDriver: true,
                }),
            ])
        ).start();
    }, []);

    const opacity = animatedValue.interpolate({
        inputRange: [0, 1],
        outputRange: [0.3, 0.7],
    });

    return (
        <Animated.View
            style={[
                {
                    width,
                    height,
                    borderRadius,
                    backgroundColor: COLORS.surfaceHighlight,
                    opacity,
                },
                style,
            ]}
        />
    );
};

export const CardSkeleton = () => (
    <View style={skeletonStyles.cardContainer}>
        <SkeletonLoader width="100%" height={400} borderRadius={SIZES.radiusL} />
        <View style={skeletonStyles.infoContainer}>
            <SkeletonLoader width={200} height={28} style={{ marginBottom: SPACING.s }} />
            <SkeletonLoader width={150} height={16} style={{ marginBottom: SPACING.m }} />
            <SkeletonLoader width="100%" height={16} style={{ marginBottom: SPACING.xs }} />
            <SkeletonLoader width="80%" height={16} style={{ marginBottom: SPACING.m }} />
            <View style={skeletonStyles.tagsContainer}>
                <SkeletonLoader width={80} height={24} borderRadius={12} />
                <SkeletonLoader width={80} height={24} borderRadius={12} />
                <SkeletonLoader width={80} height={24} borderRadius={12} />
            </View>
        </View>
    </View>
);

export const ChatSkeleton = () => (
    <View style={skeletonStyles.chatItem}>
        <SkeletonLoader width={60} height={60} borderRadius={30} />
        <View style={skeletonStyles.chatContent}>
            <SkeletonLoader width={120} height={16} style={{ marginBottom: SPACING.xs }} />
            <SkeletonLoader width="80%" height={14} />
        </View>
    </View>
);

const skeletonStyles = StyleSheet.create({
    cardContainer: {
        width: '90%',
        height: 600,
        borderRadius: SIZES.radiusL,
        backgroundColor: COLORS.surface,
        overflow: 'hidden',
    },
    infoContainer: {
        position: 'absolute',
        bottom: 0,
        left: 0,
        right: 0,
        padding: SPACING.l,
    },
    tagsContainer: {
        flexDirection: 'row',
        gap: SPACING.s,
    },
    chatItem: {
        flexDirection: 'row',
        alignItems: 'center',
        marginBottom: SPACING.l,
        paddingHorizontal: SPACING.l,
    },
    chatContent: {
        flex: 1,
        marginLeft: SPACING.m,
    },
});

