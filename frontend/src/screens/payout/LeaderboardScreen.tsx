import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, FlatList, Image, TouchableOpacity, ActivityIndicator } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { COLORS, SPACING, SIZES } from '../../theme/colors';
import { LeaderboardService } from '../../api/services';

// Fallback mock data
const MOCK_LEADERBOARD = [
    { id: '1', name: 'Selam', avatar: 'https://images.unsplash.com/photo-1627483262769-04d0a1401487?q=80&w=1000&auto=format&fit=crop', amount: 12500, rank: 1 },
    { id: '2', name: 'Tigist', avatar: 'https://images.unsplash.com/photo-1531123897727-8f129e1688ce?q=80&w=1000&auto=format&fit=crop', amount: 9800, rank: 2 },
    { id: '3', name: 'Hana', avatar: 'https://images.unsplash.com/photo-1589156280159-27698a70f29e?q=80&w=1000&auto=format&fit=crop', amount: 8750, rank: 3 },
    { id: '4', name: 'Marta', avatar: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?q=80&w=1000&auto=format&fit=crop', amount: 7200, rank: 4 },
    { id: '5', name: 'Yodit', avatar: 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?q=80&w=1000&auto=format&fit=crop', amount: 6500, rank: 5 },
    { id: '6', name: 'Beti', avatar: 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?q=80&w=1000&auto=format&fit=crop', amount: 5800, rank: 6 },
    { id: '7', name: 'Meron', avatar: 'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?q=80&w=1000&auto=format&fit=crop', amount: 5200, rank: 7 },
    { id: '8', name: 'Sara', avatar: 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?q=80&w=1000&auto=format&fit=crop', amount: 4800, rank: 8 },
];

export const LeaderboardScreen = ({ navigation }: any) => {
    const [leaderboard, setLeaderboard] = useState(MOCK_LEADERBOARD);
    const [isLoading, setIsLoading] = useState(true);
    const [timeframe, setTimeframe] = useState<'week' | 'month' | 'all'>('week');

    useEffect(() => {
        loadLeaderboard();
    }, [timeframe]);

    const loadLeaderboard = async () => {
        try {
            setIsLoading(true);
            const response = await LeaderboardService.getTopGifted(timeframe, 20);
            if (response.leaderboard && response.leaderboard.length > 0) {
                setLeaderboard(response.leaderboard);
            }
        } catch (error: any) {
            console.error('Load leaderboard error:', error);
            // Keep mock data on error
        } finally {
            setIsLoading(false);
        }
    };

    const getRankIcon = (rank: number) => {
        if (rank === 1) return '🥇';
        if (rank === 2) return '🥈';
        if (rank === 3) return '🥉';
        return `#${rank}`;
    };

    const renderLeaderboardItem = ({ item, index }: { item: any; index: number }) => {
        const isTopThree = item.rank <= 3;
        
        return (
            <View style={[styles.leaderboardItem, isTopThree && styles.leaderboardItemTop]}>
                <View style={styles.rankContainer}>
                    <Text style={styles.rankText}>{getRankIcon(item.rank)}</Text>
                </View>
                
                <Image source={{ uri: item.avatar }} style={styles.avatar} />
                
                <View style={styles.userInfo}>
                    <Text style={styles.userName}>{item.name}</Text>
                    <Text style={styles.userAmount}>🎁 {item.amount.toLocaleString()} Birr</Text>
                </View>
                
                {isTopThree && (
                    <View style={styles.badge}>
                        <Text style={styles.badgeText}>Top {item.rank}</Text>
                    </View>
                )}
            </View>
        );
    };

    return (
        <SafeAreaView style={styles.container} edges={['top']}>
            <View style={styles.header}>
                <TouchableOpacity onPress={() => navigation.goBack()} style={styles.backButton}>
                    <Text style={styles.backIcon}>←</Text>
                </TouchableOpacity>
                <Text style={styles.headerTitle}>Top Gifted 🏆</Text>
                <View style={styles.placeholder} />
            </View>

            {/* Timeframe Selector */}
            <View style={styles.timeframeContainer}>
                {(['week', 'month', 'all'] as const).map((tf) => (
                    <TouchableOpacity
                        key={tf}
                        style={[
                            styles.timeframeButton,
                            timeframe === tf && styles.timeframeButtonActive,
                        ]}
                        onPress={() => setTimeframe(tf)}
                    >
                        <Text
                            style={[
                                styles.timeframeText,
                                timeframe === tf && styles.timeframeTextActive,
                            ]}
                        >
                            {tf === 'week' ? 'This Week' : tf === 'month' ? 'This Month' : 'All Time'}
                        </Text>
                    </TouchableOpacity>
                ))}
            </View>

            {/* Top 3 Podium */}
            {isLoading ? (
                <View style={styles.loadingContainer}>
                    <ActivityIndicator size="large" color={COLORS.primary} />
                </View>
            ) : leaderboard.length >= 3 ? (
                <View style={styles.podiumContainer}>
                    <View style={styles.podiumItem}>
                        <Image
                            source={{ uri: leaderboard[1].avatar }}
                            style={[styles.podiumAvatar, styles.podiumAvatarSecond]}
                        />
                        <Text style={styles.podiumRank}>🥈</Text>
                        <Text style={styles.podiumName}>{leaderboard[1].name}</Text>
                        <Text style={styles.podiumAmount}>{leaderboard[1].amount.toLocaleString()}</Text>
                    </View>
                    
                    <View style={styles.podiumItem}>
                        <Image
                            source={{ uri: leaderboard[0].avatar }}
                            style={[styles.podiumAvatar, styles.podiumAvatarFirst]}
                        />
                        <Text style={styles.podiumRank}>🥇</Text>
                        <Text style={styles.podiumName}>{leaderboard[0].name}</Text>
                        <Text style={styles.podiumAmount}>{leaderboard[0].amount.toLocaleString()}</Text>
                    </View>
                    
                    <View style={styles.podiumItem}>
                        <Image
                            source={{ uri: leaderboard[2].avatar }}
                            style={[styles.podiumAvatar, styles.podiumAvatarThird]}
                        />
                        <Text style={styles.podiumRank}>🥉</Text>
                        <Text style={styles.podiumName}>{leaderboard[2].name}</Text>
                        <Text style={styles.podiumAmount}>{leaderboard[2].amount.toLocaleString()}</Text>
                    </View>
                </View>
            ) : null}

            {/* Leaderboard List */}
            <FlatList
                data={leaderboard.slice(3)}
                renderItem={renderLeaderboardItem}
                keyExtractor={(item) => item.id}
                contentContainerStyle={styles.listContent}
                showsVerticalScrollIndicator={false}
            />
        </SafeAreaView>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: COLORS.background,
    },
    header: {
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'space-between',
        padding: SPACING.m,
        borderBottomWidth: 1,
        borderBottomColor: COLORS.surfaceHighlight,
        backgroundColor: COLORS.surface,
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
    },
    placeholder: {
        width: 40,
    },
    timeframeContainer: {
        flexDirection: 'row',
        padding: SPACING.m,
        gap: SPACING.s,
    },
    timeframeButton: {
        flex: 1,
        padding: SPACING.s,
        borderRadius: SIZES.radiusM,
        backgroundColor: COLORS.surface,
        alignItems: 'center',
        borderWidth: 1,
        borderColor: COLORS.surfaceHighlight,
    },
    timeframeButtonActive: {
        backgroundColor: COLORS.primary,
        borderColor: COLORS.primary,
    },
    timeframeText: {
        fontSize: 12,
        color: COLORS.textSecondary,
        fontWeight: '600',
    },
    timeframeTextActive: {
        color: COLORS.background,
    },
    podiumContainer: {
        flexDirection: 'row',
        justifyContent: 'center',
        alignItems: 'flex-end',
        padding: SPACING.l,
        marginBottom: SPACING.m,
    },
    podiumItem: {
        flex: 1,
        alignItems: 'center',
    },
    podiumAvatar: {
        width: 60,
        height: 60,
        borderRadius: 30,
        borderWidth: 3,
        marginBottom: SPACING.xs,
    },
    podiumAvatarFirst: {
        width: 80,
        height: 80,
        borderRadius: 40,
        borderColor: COLORS.gold,
    },
    podiumAvatarSecond: {
        borderColor: '#C0C0C0',
    },
    podiumAvatarThird: {
        borderColor: '#CD7F32',
    },
    podiumRank: {
        fontSize: 24,
        marginBottom: SPACING.xs,
    },
    podiumName: {
        fontSize: 14,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
        marginBottom: 2,
    },
    podiumAmount: {
        fontSize: 12,
        color: COLORS.textSecondary,
    },
    listContent: {
        padding: SPACING.l,
    },
    leaderboardItem: {
        flexDirection: 'row',
        alignItems: 'center',
        backgroundColor: COLORS.surface,
        borderRadius: SIZES.radiusM,
        padding: SPACING.m,
        marginBottom: SPACING.m,
    },
    leaderboardItemTop: {
        borderWidth: 1,
        borderColor: COLORS.surfaceHighlight,
    },
    rankContainer: {
        width: 40,
        alignItems: 'center',
    },
    rankText: {
        fontSize: 18,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
    },
    avatar: {
        width: 50,
        height: 50,
        borderRadius: 25,
        marginRight: SPACING.m,
    },
    userInfo: {
        flex: 1,
    },
    userName: {
        fontSize: 16,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
        marginBottom: 2,
    },
    userAmount: {
        fontSize: 14,
        color: COLORS.textSecondary,
    },
    badge: {
        backgroundColor: COLORS.gold,
        paddingHorizontal: SPACING.s,
        paddingVertical: 4,
        borderRadius: SIZES.radiusS,
    },
    badgeText: {
        fontSize: 10,
        fontWeight: 'bold',
        color: COLORS.background,
    },
    loadingContainer: {
        padding: SPACING.xl,
        alignItems: 'center',
    },
});

