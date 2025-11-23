import React from 'react';
import { TouchableOpacity, Text, StyleSheet, Platform } from 'react-native';
import { useNavigation } from '@react-navigation/native';
import { COLORS, SPACING } from '../../theme/colors';

interface BackButtonProps {
    onPress?: () => void;
    style?: any;
}

export const BackButton: React.FC<BackButtonProps> = ({ onPress, style }) => {
    const navigation = useNavigation();

    const handlePress = () => {
        if (onPress) {
            onPress();
        } else {
            navigation.goBack();
        }
    };

    return (
        <TouchableOpacity
            style={[styles.backButton, style]}
            onPress={handlePress}
            hitSlop={{ top: 10, bottom: 10, left: 10, right: 10 }}
        >
            <Text style={styles.backIcon}>←</Text>
            <Text style={styles.backText}>Back</Text>
        </TouchableOpacity>
    );
};

const styles = StyleSheet.create({
    backButton: {
        flexDirection: 'row',
        alignItems: 'center',
        paddingVertical: SPACING.s,
        paddingHorizontal: SPACING.m,
        marginBottom: SPACING.m,
    },
    backIcon: {
        fontSize: 24,
        color: COLORS.textSecondary,
        marginRight: SPACING.xs,
    },
    backText: {
        fontSize: 16,
        color: COLORS.textSecondary,
        fontWeight: '500',
    },
});
