import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity, Dimensions } from 'react-native';
import { COLORS, SPACING, SIZES } from '../../theme/colors';
import { RelationshipGoal } from '../../constants/ethiopianData';

interface RelationshipGoalCardProps {
    goal: RelationshipGoal;
    isSelected: boolean;
    onPress: () => void;
}

const { width } = Dimensions.get('window');
const cardWidth = (width - SPACING.l * 3) / 2; // 2 columns with spacing

export const RelationshipGoalCard: React.FC<RelationshipGoalCardProps> = ({
    goal,
    isSelected,
    onPress,
}) => {
    return (
        <TouchableOpacity
            style={[
                styles.card,
                isSelected && styles.cardSelected,
            ]}
            onPress={onPress}
            activeOpacity={0.7}
        >
            <Text style={styles.emoji}>{goal.emoji}</Text>
            <Text style={[styles.title, isSelected && styles.titleSelected]}>
                {goal.title}
            </Text>
            <Text style={styles.subtitle}>{goal.subtitle}</Text>

            {isSelected && (
                <View style={styles.checkmark}>
                    <Text style={styles.checkmarkText}>✓</Text>
                </View>
            )}
        </TouchableOpacity>
    );
};

const styles = StyleSheet.create({
    card: {
        width: cardWidth,
        backgroundColor: COLORS.surface,
        borderRadius: SIZES.radiusL,
        padding: SPACING.m,
        marginBottom: SPACING.m,
        borderWidth: 2,
        borderColor: COLORS.surfaceHighlight,
        alignItems: 'center',
        minHeight: 140,
        justifyContent: 'center',
    },
    cardSelected: {
        borderColor: COLORS.primary,
        backgroundColor: 'rgba(167, 255, 131, 0.1)',
    },
    emoji: {
        fontSize: 40,
        marginBottom: SPACING.s,
    },
    title: {
        fontSize: 14,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
        textAlign: 'center',
        marginBottom: 4,
    },
    titleSelected: {
        color: COLORS.primary,
    },
    subtitle: {
        fontSize: 11,
        color: COLORS.textSecondary,
        textAlign: 'center',
    },
    checkmark: {
        position: 'absolute',
        top: 8,
        right: 8,
        width: 24,
        height: 24,
        borderRadius: 12,
        backgroundColor: COLORS.primary,
        alignItems: 'center',
        justifyContent: 'center',
    },
    checkmarkText: {
        color: COLORS.background,
        fontSize: 14,
        fontWeight: 'bold',
    },
});
