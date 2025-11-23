import React from 'react';
import { TouchableOpacity, Text, StyleSheet, ActivityIndicator, ViewStyle, TextStyle } from 'react-native';
import { COLORS, SIZES, SPACING } from '../../theme/colors';

interface ButtonProps {
    title: string;
    onPress: () => void;
    variant?: 'primary' | 'secondary' | 'outline' | 'ghost';
    size?: 'small' | 'medium' | 'large';
    isLoading?: boolean;
    disabled?: boolean;
    style?: ViewStyle;
    textStyle?: TextStyle;
    icon?: React.ReactNode;
}

export const Button: React.FC<ButtonProps> = ({
    title,
    onPress,
    variant = 'primary',
    size = 'medium',
    isLoading = false,
    disabled = false,
    style,
    textStyle,
    icon,
}) => {
    const getBackgroundColor = () => {
        if (disabled) return COLORS.surfaceHighlight;
        switch (variant) {
            case 'primary': return COLORS.primary;
            case 'secondary': return COLORS.surfaceHighlight;
            case 'outline': return 'transparent';
            case 'ghost': return 'transparent';
            default: return COLORS.primary;
        }
    };

    const getTextColor = () => {
        if (disabled) return COLORS.textSecondary;
        switch (variant) {
            case 'primary': return COLORS.background; // Black text on Lime
            case 'secondary': return COLORS.textPrimary;
            case 'outline': return COLORS.primary;
            case 'ghost': return COLORS.textSecondary;
            default: return COLORS.background;
        }
    };

    const getHeight = () => {
        switch (size) {
            case 'small': return 36;
            case 'medium': return 48;
            case 'large': return 56;
            default: return 48;
        }
    };

    return (
        <TouchableOpacity
            onPress={onPress}
            disabled={disabled || isLoading}
            style={[
                styles.container,
                {
                    backgroundColor: getBackgroundColor(),
                    height: getHeight(),
                    borderColor: variant === 'outline' ? COLORS.primary : 'transparent',
                    borderWidth: variant === 'outline' ? 1 : 0,
                },
                style,
            ]}
        >
            {isLoading ? (
                <ActivityIndicator color={variant === 'primary' ? COLORS.background : COLORS.primary} />
            ) : (
                <>
                    {icon && icon}
                    <Text
                        style={[
                            styles.text,
                            {
                                color: getTextColor(),
                                fontSize: size === 'small' ? 14 : 16,
                                marginLeft: icon ? SPACING.s : 0,
                            },
                            textStyle,
                        ]}
                    >
                        {title}
                    </Text>
                </>
            )}
        </TouchableOpacity>
    );
};

const styles = StyleSheet.create({
    container: {
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'center',
        borderRadius: SIZES.radiusM,
        paddingHorizontal: SPACING.l,
    },
    text: {
        fontWeight: 'bold',
        letterSpacing: 0.5,
    },
});
