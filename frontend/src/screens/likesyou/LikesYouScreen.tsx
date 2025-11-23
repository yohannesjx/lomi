import React, { useState, useEffect } from 'react';
import {
    View,
    Text,
    StyleSheet,
    ScrollView,
    TouchableOpacity,
    Image,
    ActivityIndicator,
    Alert,
    Animated,
    Dimensions,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { COLORS, SPACING, SIZES } from '../../theme/colors';
import { LikesService } from '../../api/services';
import { CoinService } from '../../api/services';
import { useAuthStore } from '../../store/authStore';
// Simple confetti effect - can be replaced with react-native-confetti-cannon if installed

const { width } = Dimensions.get('window');

interface PendingLike {
    user: {
        id: string;
        name: string;
        city: string;
        avatar?: string;
    };
    liked_at: string;
    is_revealed: boolean;
}

interface PendingLikesResponse {
    pending_likes: PendingLike[];
    count: number;
    has_free_reveal: boolean;
    reset_at: string;
}

// BlurCard Component with glowing silhouette
const BlurCard = ({ 
    like, 
    isRevealed, 
    onReveal 
}: { 
    like: PendingLike; 
    isRevealed: boolean;
    onReveal: () => void;
}) => {
    const glowAnim = React.useRef(new Animated.Value(0)).current;
    const [revealed, setRevealed] = useState(isRevealed);

    useEffect(() => {
        if (!revealed) {
            Animated.loop(
                Animated.sequence([
                    Animated.timing(glowAnim, {
                        toValue: 1,
                        duration: 2000,
                        useNativeDriver: true,
                    }),
                    Animated.timing(glowAnim, {
                        toValue: 0,
                        duration: 2000,
                        useNativeDriver: true,
                    }),
                ])
            ).start();
        }
    }, [revealed]);

    const glowOpacity = glowAnim.interpolate({
        inputRange: [0, 1],
        outputRange: [0.3, 0.7],
    });

    const firstName = like.user.name.split(' ')[0];
    const maskedName = `${firstName[0]}***`;

    const timeAgo = (dateString: string) => {
        const date = new Date(dateString);
        const now = new Date();
        const diffMs = now.getTime() - date.getTime();
        const diffMins = Math.floor(diffMs / 60000);
        const diffHours = Math.floor(diffMins / 60);
        const diffDays = Math.floor(diffHours / 24);

        if (diffMins < 60) return `${diffMins}m ago`;
        if (diffHours < 24) return `${diffHours}h ago`;
        return `${diffDays}d ago`;
    };

    return (
        <TouchableOpacity
            style={styles.blurCard}
            onPress={!revealed ? onReveal : undefined}
            activeOpacity={0.8}
        >
            <View style={styles.blurCardContent}>
                {!revealed ? (
                    <>
                        {/* Blurred silhouette with glow */}
                        <Animated.View
                            style={[
                                styles.blurAvatarContainer,
                                { opacity: glowOpacity },
                            ]}
                        >
                            <View style={styles.blurAvatar}>
                                <View style={styles.blurOverlay} />
                                <View style={styles.silhouetteGlow} />
                            </View>
                        </Animated.View>
                        <Text style={styles.blurName}>{maskedName}</Text>
                        <Text style={styles.blurCity}>From {like.user.city}</Text>
                        <Text style={styles.blurTime}>Liked you {timeAgo(like.liked_at)}</Text>
                        <View style={styles.revealHint}>
                            <Text style={styles.revealHintText}>Tap to reveal 👆</Text>
                        </View>
                    </>
                ) : (
                    <>
                        {/* Revealed profile */}
                        {like.user.avatar ? (
                            <Image
                                source={{ uri: like.user.avatar }}
                                style={styles.revealedAvatar}
                            />
                        ) : (
                            <View style={styles.revealedAvatarPlaceholder}>
                                <Text style={styles.revealedAvatarText}>
                                    {firstName[0].toUpperCase()}
                                </Text>
                            </View>
                        )}
                        <Text style={styles.revealedName}>{like.user.name}</Text>
                        <Text style={styles.revealedCity}>From {like.user.city}</Text>
                        <Text style={styles.revealedTime}>Liked you {timeAgo(like.liked_at)}</Text>
                    </>
                )}
            </View>
        </TouchableOpacity>
    );
};

export const LikesYouScreen = ({ navigation }: any) => {
    const { user } = useAuthStore();
    const [pendingLikes, setPendingLikes] = useState<PendingLike[]>([]);
    const [count, setCount] = useState(0);
    const [hasFreeReveal, setHasFreeReveal] = useState(false);
    const [resetAt, setResetAt] = useState<Date | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [isRevealing, setIsRevealing] = useState(false);
    const [coinBalance, setCoinBalance] = useState(0);
    const [language, setLanguage] = useState<'en' | 'am'>('en');
    const [showConfetti, setShowConfetti] = useState(false);
    const [countdown, setCountdown] = useState('');

    useEffect(() => {
        loadPendingLikes();
        loadCoinBalance();
        startCountdown();
    }, []);

    useEffect(() => {
        if (resetAt) {
            const interval = setInterval(() => {
                updateCountdown();
            }, 1000);
            return () => clearInterval(interval);
        }
    }, [resetAt]);

    const loadPendingLikes = async () => {
        try {
            setIsLoading(true);
            const response: PendingLikesResponse = await LikesService.getPendingLikes();
            setPendingLikes(response.pending_likes || []);
            setCount(response.count || 0);
            setHasFreeReveal(response.has_free_reveal || false);
            if (response.reset_at) {
                setResetAt(new Date(response.reset_at));
            }
        } catch (error: any) {
            console.error('Load pending likes error:', error);
            if (error.response?.status === 401) {
                // User not authenticated - show message
                if (__DEV__) {
                    console.warn('⚠️ Not authenticated. Please log in first.');
                }
                Alert.alert(
                    'Authentication Required',
                    'Please log in to see who likes you.',
                    [
                        {
                            text: 'OK',
                            onPress: () => navigation.navigate('Welcome'),
                        },
                    ]
                );
            } else {
                Alert.alert('Error', 'Failed to load likes');
            }
        } finally {
            setIsLoading(false);
        }
    };

    const loadCoinBalance = async () => {
        try {
            const response = await CoinService.getBalance();
            setCoinBalance(response.coin_balance || 0);
        } catch (error: any) {
            console.error('Load coin balance error:', error);
            // In dev mode without auth, this is expected - set balance to 0
            if (error.response?.status === 401 && __DEV__) {
                setCoinBalance(0);
            }
        }
    };

    const startCountdown = () => {
        if (resetAt) {
            updateCountdown();
        }
    };

    const updateCountdown = () => {
        if (!resetAt) return;
        const now = new Date();
        const diff = resetAt.getTime() - now.getTime();

        if (diff <= 0) {
            setCountdown('');
            loadPendingLikes(); // Reload to get new free reveal
            return;
        }

        const hours = Math.floor(diff / (1000 * 60 * 60));
        const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
        setCountdown(`${hours}h ${minutes}m`);
    };

    const handleRevealAll = async () => {
        if (isRevealing) return;

        // Check if user has coins or free reveal
        if (!hasFreeReveal && coinBalance < 299) {
            // Redirect to coin purchase
            navigation.navigate('BuyCoins', { 
                preselectedPackage: 500,
                message: 'Unlock your secret admirers now!'
            });
            return;
        }

        setIsRevealing(true);
        try {
            const response = await LikesService.revealLike({ reveal_all: true });
            
            if (response.revealed_users && response.revealed_users.length > 0) {
                setShowConfetti(true);
                setTimeout(() => setShowConfetti(false), 3000);
                
                // Update local state
                const updatedLikes = pendingLikes.map(like => {
                    const revealed = response.revealed_users.find((u: any) => u.id === like.user.id);
                    return revealed ? { ...like, is_revealed: true, user: revealed } : like;
                });
                setPendingLikes(updatedLikes);
                setCount(0);
                setHasFreeReveal(false);
                setCoinBalance(response.new_balance || 0);

                Alert.alert(
                    language === 'am' ? 'Enesu fikirwo yeshalewal!' : 'All revealed!',
                    language === 'am' 
                        ? 'Now go say hi 😉' 
                        : 'All your admirers are revealed! Now go say hi 😉',
                    [{ text: 'OK' }]
                );
            }
        } catch (error: any) {
            if (error.response?.status === 400 && error.response?.data?.error === 'Insufficient coins') {
                navigation.navigate('BuyCoins', { 
                    preselectedPackage: 500,
                    message: 'Unlock your secret admirers now!'
                });
            } else {
                Alert.alert('Error', error.response?.data?.error || 'Failed to reveal');
            }
        } finally {
            setIsRevealing(false);
        }
    };

    const handleRevealOne = async (targetId: string) => {
        if (isRevealing) return;

        // Check if user has coins or free reveal
        if (!hasFreeReveal && coinBalance < 99) {
            navigation.navigate('BuyCoins', { 
                preselectedPackage: 500,
                message: 'Unlock your secret admirers now!'
            });
            return;
        }

        setIsRevealing(true);
        try {
            const response = await LikesService.revealLike({ 
                reveal_all: false,
                target_id: targetId 
            });
            
            if (response.revealed_users && response.revealed_users.length > 0) {
                const revealedUser = response.revealed_users[0];
                
                // Update local state
                const updatedLikes = pendingLikes.map(like => 
                    like.user.id === targetId 
                        ? { ...like, is_revealed: true, user: revealedUser }
                        : like
                );
                setPendingLikes(updatedLikes);
                setCount(count - 1);
                setHasFreeReveal(false);
                setCoinBalance(response.new_balance || 0);
            }
        } catch (error: any) {
            if (error.response?.status === 400 && error.response?.data?.error === 'Insufficient coins') {
                navigation.navigate('BuyCoins', { 
                    preselectedPackage: 500,
                    message: 'Unlock your secret admirers now!'
                });
            } else {
                Alert.alert('Error', error.response?.data?.error || 'Failed to reveal');
            }
        } finally {
            setIsRevealing(false);
        }
    };

    // If user has 0 coins and no free reveal, redirect to coin purchase
    useEffect(() => {
        if (!isLoading && count > 0 && !hasFreeReveal && coinBalance === 0) {
            // Don't auto-redirect, let user see the screen first
        }
    }, [isLoading, count, hasFreeReveal, coinBalance]);

    if (isLoading) {
        return (
            <SafeAreaView style={styles.container} edges={['top']}>
                <View style={styles.loadingContainer}>
                    <ActivityIndicator size="large" color={COLORS.primary} />
                    <Text style={styles.loadingText}>Loading your admirers...</Text>
                </View>
            </SafeAreaView>
        );
    }

    if (count === 0) {
        return (
            <SafeAreaView style={styles.container} edges={['top']}>
                <View style={styles.emptyContainer}>
                    <Text style={styles.emptyIcon}>💔</Text>
                    <Text style={styles.emptyTitle}>
                        {language === 'am' ? 'Yelew fikir yelew' : 'No likes yet'}
                    </Text>
                    <Text style={styles.emptySubtitle}>
                        {language === 'am' 
                            ? 'Keep swiping to get more likes!' 
                            : 'Keep swiping to get more likes!'}
                    </Text>
                    <TouchableOpacity
                        style={styles.discoverButton}
                        onPress={() => navigation.navigate('Discover')}
                    >
                        <Text style={styles.discoverButtonText}>Start Discovering 🔥</Text>
                    </TouchableOpacity>
                </View>
            </SafeAreaView>
        );
    }

    return (
        <SafeAreaView style={styles.container} edges={['top']}>
            {showConfetti && (
                <View style={styles.confettiContainer} pointerEvents="none">
                    <Text style={styles.confettiText}>🎉✨🎊</Text>
                </View>
            )}

            {/* Header */}
            <View style={styles.header}>
                <TouchableOpacity onPress={() => navigation.goBack()} style={styles.backButton}>
                    <Text style={styles.backIcon}>←</Text>
                </TouchableOpacity>
                <Text style={styles.headerTitle}>
                    {language === 'am' ? 'Ye Fikir List' : 'Who Likes You'}
                </Text>
                <TouchableOpacity onPress={() => setLanguage(language === 'en' ? 'am' : 'en')}>
                    <Text style={styles.languageToggle}>
                        {language === 'en' ? 'አማ' : 'EN'}
                    </Text>
                </TouchableOpacity>
            </View>

            <ScrollView style={styles.content} showsVerticalScrollIndicator={false}>
                {/* Big Header */}
                <View style={styles.bigHeader}>
                    <Text style={styles.bigHeaderText}>
                        {language === 'am' 
                            ? `${count} nefsoch fikirwo yeshalewal! 😍`
                            : `${count} people liked you! 😍`}
                    </Text>
                    <Text style={styles.bigHeaderSubtext}>
                        {language === 'am'
                            ? 'Unlock to see who they are'
                            : 'Unlock to see who they are'}
                    </Text>
                </View>

                {/* Free Reveal Banner */}
                {hasFreeReveal && (
                    <View style={styles.freeBanner}>
                        <Text style={styles.freeBannerText}>
                            🎁 {language === 'am' ? '1 FREE reveal today!' : '1 FREE reveal today!'}
                        </Text>
                        {countdown && (
                            <Text style={styles.freeBannerCountdown}>
                                {language === 'am' ? 'Resets in' : 'Resets in'} {countdown}
                            </Text>
                        )}
                    </View>
                )}

                {/* Reveal All Button */}
                <TouchableOpacity
                    style={[styles.revealAllButton, isRevealing && styles.revealAllButtonDisabled]}
                    onPress={handleRevealAll}
                    disabled={isRevealing}
                >
                    {isRevealing ? (
                        <ActivityIndicator color={COLORS.background} />
                    ) : (
                        <>
                            <Text style={styles.revealAllButtonText}>
                                {language === 'am'
                                    ? `Reveal all ${count} for ${hasFreeReveal ? 'FREE' : '299 coins'}`
                                    : `Reveal all ${count} for ${hasFreeReveal ? 'FREE' : '299 coins'}`}
                            </Text>
                            {!hasFreeReveal && (
                                <Text style={styles.revealAllButtonSubtext}>
                                    {language === 'am' ? 'Recommended' : 'Recommended'} ⭐
                                </Text>
                            )}
                        </>
                    )}
                </TouchableOpacity>

                {/* Reveal One Button */}
                <TouchableOpacity
                    style={styles.revealOneButton}
                    onPress={() => {
                        if (pendingLikes.length > 0 && !pendingLikes[0].is_revealed) {
                            handleRevealOne(pendingLikes[0].user.id);
                        }
                    }}
                    disabled={isRevealing || (pendingLikes.length > 0 && pendingLikes[0].is_revealed)}
                >
                    <Text style={styles.revealOneButtonText}>
                        {language === 'am'
                            ? `Reveal one by one – ${hasFreeReveal ? 'FREE' : '99 coins'} each`
                            : `Reveal one by one – ${hasFreeReveal ? 'FREE' : '99 coins'} each`}
                    </Text>
                </TouchableOpacity>

                {/* Likes Grid */}
                <View style={styles.likesGrid}>
                    {pendingLikes.map((like, index) => (
                        <BlurCard
                            key={like.user.id}
                            like={like}
                            isRevealed={like.is_revealed}
                            onReveal={() => handleRevealOne(like.user.id)}
                        />
                    ))}
                </View>

                {/* Countdown Timer */}
                {countdown && (
                    <View style={styles.countdownContainer}>
                        <Text style={styles.countdownText}>
                            {language === 'am'
                                ? `Your likes reset in ${countdown}`
                                : `Your likes reset in ${countdown}`}
                        </Text>
                    </View>
                )}
            </ScrollView>
        </SafeAreaView>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: COLORS.background,
    },
    loadingContainer: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
    },
    loadingText: {
        color: COLORS.textSecondary,
        marginTop: SPACING.m,
        fontSize: 16,
    },
    header: {
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'space-between',
        padding: SPACING.m,
        borderBottomWidth: 1,
        borderBottomColor: COLORS.surfaceHighlight,
    },
    backButton: {
        padding: SPACING.s,
    },
    backIcon: {
        fontSize: 24,
        color: COLORS.textPrimary,
    },
    headerTitle: {
        fontSize: 20,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
        flex: 1,
        textAlign: 'center',
    },
    languageToggle: {
        color: COLORS.primary,
        fontSize: 14,
        fontWeight: '600',
        padding: SPACING.s,
    },
    content: {
        flex: 1,
    },
    bigHeader: {
        padding: SPACING.xl,
        alignItems: 'center',
    },
    bigHeaderText: {
        fontSize: 28,
        fontWeight: 'bold',
        color: COLORS.primary,
        textAlign: 'center',
        marginBottom: SPACING.s,
    },
    bigHeaderSubtext: {
        fontSize: 16,
        color: COLORS.textSecondary,
        textAlign: 'center',
    },
    freeBanner: {
        backgroundColor: 'rgba(255, 215, 0, 0.15)',
        borderWidth: 1,
        borderColor: COLORS.gold,
        borderRadius: SIZES.radiusM,
        padding: SPACING.m,
        marginHorizontal: SPACING.l,
        marginBottom: SPACING.m,
        alignItems: 'center',
    },
    freeBannerText: {
        color: COLORS.gold,
        fontSize: 18,
        fontWeight: 'bold',
        marginBottom: SPACING.xs,
    },
    freeBannerCountdown: {
        color: COLORS.gold,
        fontSize: 14,
    },
    revealAllButton: {
        backgroundColor: COLORS.primary,
        borderRadius: SIZES.radiusL,
        padding: SPACING.l,
        marginHorizontal: SPACING.l,
        marginBottom: SPACING.m,
        alignItems: 'center',
        shadowColor: COLORS.primary,
        shadowOffset: { width: 0, height: 0 },
        shadowOpacity: 0.5,
        shadowRadius: 20,
        elevation: 10,
    },
    revealAllButtonDisabled: {
        opacity: 0.6,
    },
    revealAllButtonText: {
        color: COLORS.background,
        fontSize: 20,
        fontWeight: 'bold',
    },
    revealAllButtonSubtext: {
        color: COLORS.background,
        fontSize: 14,
        marginTop: SPACING.xs,
        opacity: 0.9,
    },
    revealOneButton: {
        backgroundColor: COLORS.surface,
        borderWidth: 1,
        borderColor: COLORS.primary,
        borderRadius: SIZES.radiusM,
        padding: SPACING.m,
        marginHorizontal: SPACING.l,
        marginBottom: SPACING.l,
        alignItems: 'center',
    },
    revealOneButtonText: {
        color: COLORS.primary,
        fontSize: 16,
        fontWeight: '600',
    },
    likesGrid: {
        flexDirection: 'row',
        flexWrap: 'wrap',
        justifyContent: 'space-between',
        paddingHorizontal: SPACING.l,
        paddingBottom: SPACING.xl,
    },
    blurCard: {
        width: (width - SPACING.l * 3) / 2,
        backgroundColor: COLORS.surface,
        borderRadius: SIZES.radiusM,
        padding: SPACING.m,
        marginBottom: SPACING.m,
        borderWidth: 1,
        borderColor: COLORS.surfaceHighlight,
    },
    blurCardContent: {
        alignItems: 'center',
    },
    blurAvatarContainer: {
        position: 'relative',
        marginBottom: SPACING.m,
    },
    blurAvatar: {
        width: 80,
        height: 80,
        borderRadius: 40,
        backgroundColor: COLORS.surfaceHighlight,
        overflow: 'hidden',
        position: 'relative',
    },
    blurOverlay: {
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        backgroundColor: 'rgba(0, 0, 0, 0.7)',
        backdropFilter: 'blur(10px)',
    },
    silhouetteGlow: {
        position: 'absolute',
        top: -10,
        left: -10,
        right: -10,
        bottom: -10,
        borderRadius: 50,
        backgroundColor: COLORS.primary,
        opacity: 0.3,
    },
    blurName: {
        fontSize: 18,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
        marginBottom: SPACING.xs,
    },
    blurCity: {
        fontSize: 14,
        color: COLORS.textSecondary,
        marginBottom: SPACING.xs,
    },
    blurTime: {
        fontSize: 12,
        color: COLORS.primary,
        marginBottom: SPACING.s,
    },
    revealHint: {
        backgroundColor: 'rgba(167, 255, 131, 0.1)',
        paddingHorizontal: SPACING.s,
        paddingVertical: SPACING.xs,
        borderRadius: SIZES.radiusS,
        marginTop: SPACING.xs,
    },
    revealHintText: {
        color: COLORS.primary,
        fontSize: 12,
    },
    revealedAvatar: {
        width: 80,
        height: 80,
        borderRadius: 40,
        marginBottom: SPACING.m,
    },
    revealedAvatarPlaceholder: {
        width: 80,
        height: 80,
        borderRadius: 40,
        backgroundColor: COLORS.primary,
        justifyContent: 'center',
        alignItems: 'center',
        marginBottom: SPACING.m,
    },
    revealedAvatarText: {
        fontSize: 32,
        fontWeight: 'bold',
        color: COLORS.background,
    },
    revealedName: {
        fontSize: 18,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
        marginBottom: SPACING.xs,
    },
    revealedCity: {
        fontSize: 14,
        color: COLORS.textSecondary,
        marginBottom: SPACING.xs,
    },
    revealedTime: {
        fontSize: 12,
        color: COLORS.primary,
    },
    countdownContainer: {
        padding: SPACING.l,
        alignItems: 'center',
        marginBottom: SPACING.xl,
    },
    countdownText: {
        color: COLORS.textSecondary,
        fontSize: 14,
    },
    emptyContainer: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        padding: SPACING.xl,
    },
    emptyIcon: {
        fontSize: 64,
        marginBottom: SPACING.l,
    },
    emptyTitle: {
        fontSize: 24,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
        marginBottom: SPACING.s,
        textAlign: 'center',
    },
    emptySubtitle: {
        fontSize: 16,
        color: COLORS.textSecondary,
        textAlign: 'center',
        marginBottom: SPACING.xl,
    },
    discoverButton: {
        backgroundColor: COLORS.primary,
        paddingHorizontal: SPACING.xl,
        paddingVertical: SPACING.m,
        borderRadius: SIZES.radiusM,
    },
    discoverButtonText: {
        color: COLORS.background,
        fontSize: 18,
        fontWeight: 'bold',
    },
    confettiContainer: {
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        justifyContent: 'center',
        alignItems: 'center',
        zIndex: 1000,
    },
    confettiText: {
        fontSize: 64,
    },
});

