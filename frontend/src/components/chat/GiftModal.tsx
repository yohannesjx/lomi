import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, Modal, TouchableOpacity, FlatList, Image, Dimensions, ActivityIndicator, Alert } from 'react-native';
import { BlurView } from 'expo-blur';
import { COLORS, SPACING, SIZES } from '../../theme/colors';
import { GiftService } from '../../api/services';
import { GiftAnimation } from './GiftAnimation';

const { width } = Dimensions.get('window');
const COLUMN_COUNT = 3;
const ITEM_WIDTH = (width - (SPACING.l * 2) - (SPACING.m * (COLUMN_COUNT - 1))) / COLUMN_COUNT;

interface GiftModalProps {
    visible: boolean;
    onClose: () => void;
    onSendGift: (gift: any) => void;
    coinBalance: number;
    onBuyCoins?: () => void;
}

export const GiftModal: React.FC<GiftModalProps> = ({ visible, onClose, onSendGift, coinBalance, onBuyCoins }) => {
    const [gifts, setGifts] = useState<any[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [showAnimation, setShowAnimation] = useState(false);
    const [selectedGift, setSelectedGift] = useState<any>(null);

    useEffect(() => {
        if (visible) {
            loadGifts();
        }
    }, [visible]);

    const loadGifts = async () => {
        try {
            setIsLoading(true);
            const response = await GiftService.getGifts();
            setGifts(response.gifts || []);
        } catch (error: any) {
            console.error('Load gifts error:', error);
            Alert.alert('Error', 'Failed to load gifts');
        } finally {
            setIsLoading(false);
        }
    };

    const handleGiftSelect = (gift: any) => {
        if (coinBalance < gift.coin_price) {
            Alert.alert(
                'Insufficient Coins',
                `You need ${gift.coin_price} coins to send this gift. You have ${coinBalance} coins.`,
                [
                    { text: 'Cancel', style: 'cancel' },
                    {
                        text: 'Buy Coins',
                        onPress: () => {
                            onClose();
                            onBuyCoins?.();
                        },
                    },
                ]
            );
            return;
        }

        setSelectedGift(gift);
        setShowAnimation(true);
    };

    const handleAnimationComplete = () => {
        setShowAnimation(false);
        onSendGift(selectedGift);
        setSelectedGift(null);
    };

    const renderGiftItem = ({ item }: { item: any }) => {
        const canAfford = coinBalance >= item.coin_price;
        return (
            <TouchableOpacity
                style={[
                    styles.giftItem,
                    !canAfford && styles.giftItemDisabled,
                    item.is_featured && styles.giftItemFeatured,
                ]}
                onPress={() => handleGiftSelect(item)}
                disabled={!canAfford}
            >
                {item.is_featured && (
                    <View style={styles.featuredBadge}>
                        <Text style={styles.featuredText}>⭐</Text>
                    </View>
                )}
                <View style={styles.giftIconContainer}>
                    {item.icon_url ? (
                        <Image source={{ uri: item.icon_url }} style={styles.giftIconImage} />
                    ) : (
                        <Text style={styles.giftIcon}>🎁</Text>
                    )}
                </View>
                <Text style={styles.giftName} numberOfLines={1}>
                    {item.name_en || item.name_am || 'Gift'}
                </Text>
                <View style={styles.priceTag}>
                    <Text style={[styles.priceText, !canAfford && styles.priceTextDisabled]}>
                        💎 {item.coin_price}
                    </Text>
                </View>
                {!canAfford && (
                    <Text style={styles.insufficientText}>Need {item.coin_price - coinBalance} more</Text>
                )}
            </TouchableOpacity>
        );
    };

    return (
        <Modal
            visible={visible}
            transparent
            animationType="slide"
            onRequestClose={onClose}
        >
            <View style={styles.overlay}>
                <TouchableOpacity style={styles.backdrop} onPress={onClose} />

                <View style={styles.modalContent}>
                    <View style={styles.header}>
                        <Text style={styles.title}>Send a Gift 🎁</Text>
                        <View style={styles.balanceContainer}>
                            <Text style={styles.balanceText}>Balance: 💎 {coinBalance}</Text>
                        </View>
                    </View>

                    {isLoading ? (
                        <View style={styles.loadingContainer}>
                            <ActivityIndicator size="large" color={COLORS.primary} />
                            <Text style={styles.loadingText}>Loading gifts...</Text>
                        </View>
                    ) : gifts.length === 0 ? (
                        <View style={styles.emptyContainer}>
                            <Text style={styles.emptyText}>No gifts available</Text>
                        </View>
                    ) : (
                        <FlatList
                            data={gifts}
                            renderItem={renderGiftItem}
                            keyExtractor={item => item.id}
                            numColumns={COLUMN_COUNT}
                            columnWrapperStyle={styles.row}
                            contentContainerStyle={styles.listContent}
                        />
                    )}

                    <TouchableOpacity
                        style={styles.buyCoinsButton}
                        onPress={() => {
                            onClose();
                            onBuyCoins?.();
                        }}
                    >
                        <Text style={styles.buyCoinsText}>+ Buy More Coins</Text>
                    </TouchableOpacity>

                    <GiftAnimation
                        visible={showAnimation}
                        gift={selectedGift || {}}
                        onComplete={handleAnimationComplete}
                    />
                </View>
            </View>
        </Modal>
    );
};

const styles = StyleSheet.create({
    overlay: {
        flex: 1,
        justifyContent: 'flex-end',
    },
    backdrop: {
        ...StyleSheet.absoluteFillObject,
        backgroundColor: 'rgba(0,0,0,0.5)',
    },
    modalContent: {
        backgroundColor: COLORS.surface,
        borderTopLeftRadius: SIZES.radiusXL,
        borderTopRightRadius: SIZES.radiusXL,
        height: '60%',
        padding: SPACING.l,
    },
    header: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
        marginBottom: SPACING.l,
    },
    title: {
        fontSize: 20,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
    },
    balanceContainer: {
        backgroundColor: 'rgba(255, 215, 0, 0.1)',
        paddingHorizontal: SPACING.m,
        paddingVertical: SPACING.xs,
        borderRadius: SIZES.radiusM,
        borderWidth: 1,
        borderColor: COLORS.gold,
    },
    balanceText: {
        color: COLORS.gold,
        fontWeight: 'bold',
        fontSize: 14,
    },
    listContent: {
        paddingBottom: SPACING.xl,
    },
    row: {
        gap: SPACING.m,
        marginBottom: SPACING.m,
    },
    giftItem: {
        width: ITEM_WIDTH,
        alignItems: 'center',
        backgroundColor: COLORS.background,
        padding: SPACING.s,
        borderRadius: SIZES.radiusM,
        borderWidth: 1,
        borderColor: COLORS.surfaceHighlight,
    },
    giftIconContainer: {
        width: 60,
        height: 60,
        alignItems: 'center',
        justifyContent: 'center',
        marginBottom: SPACING.xs,
    },
    giftIcon: {
        fontSize: 40,
    },
    giftName: {
        color: COLORS.textPrimary,
        fontSize: 14,
        fontWeight: '500',
        marginBottom: 4,
    },
    priceTag: {
        backgroundColor: COLORS.surfaceHighlight,
        paddingHorizontal: SPACING.s,
        paddingVertical: 2,
        borderRadius: SIZES.radiusS,
    },
    priceText: {
        color: COLORS.gold,
        fontSize: 12,
        fontWeight: 'bold',
    },
    buyCoinsButton: {
        backgroundColor: COLORS.primary,
        padding: SPACING.m,
        borderRadius: SIZES.radiusM,
        alignItems: 'center',
        marginTop: 'auto',
    },
    buyCoinsText: {
        color: COLORS.background,
        fontWeight: 'bold',
        fontSize: 16,
    },
    giftItemDisabled: {
        opacity: 0.5,
    },
    giftItemFeatured: {
        borderColor: COLORS.gold,
        borderWidth: 2,
    },
    featuredBadge: {
        position: 'absolute',
        top: -8,
        right: -8,
        backgroundColor: COLORS.gold,
        borderRadius: 12,
        width: 24,
        height: 24,
        alignItems: 'center',
        justifyContent: 'center',
        zIndex: 1,
    },
    featuredText: {
        fontSize: 14,
    },
    giftIconImage: {
        width: 60,
        height: 60,
        borderRadius: 30,
    },
    insufficientText: {
        fontSize: 10,
        color: COLORS.textSecondary,
        marginTop: 4,
        textAlign: 'center',
    },
    priceTextDisabled: {
        color: COLORS.textSecondary,
    },
    loadingContainer: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        paddingVertical: SPACING.xl,
    },
    loadingText: {
        marginTop: SPACING.m,
        color: COLORS.textSecondary,
    },
    emptyContainer: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        paddingVertical: SPACING.xl,
    },
    emptyText: {
        color: COLORS.textSecondary,
        fontSize: 16,
    },
});
