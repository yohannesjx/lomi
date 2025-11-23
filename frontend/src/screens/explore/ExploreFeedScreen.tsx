import React, { useState, useEffect, useRef } from 'react';
import {
    View,
    Text,
    StyleSheet,
    FlatList,
    Image,
    TouchableOpacity,
    Dimensions,
    RefreshControl,
    ActivityIndicator,
    ScrollView,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { COLORS, SPACING, SIZES } from '../../theme/colors';
import { DiscoveryService } from '../../api/services';
import { LinearGradient } from 'expo-linear-gradient';

const { width } = Dimensions.get('window');
const NUM_COLUMNS = 2;
const GAP = 2; // Gap between items
const DIVIDER_WIDTH = 2; // Divider line width
const ITEM_WIDTH = (width - GAP - (NUM_COLUMNS - 1) * DIVIDER_WIDTH) / NUM_COLUMNS;

// Mock explore feed data with various image heights
const MOCK_FEED = [
    { id: '1', url: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?q=80&w=800&auto=format&fit=crop', height: 300, type: 'photo', user: { name: 'Dawit', avatar: 'https://images.unsplash.com/photo-1531384441138-2736e62e0919?q=80&w=100&auto=format&fit=crop' } },
    { id: '2', url: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?q=80&w=800&auto=format&fit=crop', height: 450, type: 'photo', user: { name: 'Selam', avatar: 'https://images.unsplash.com/photo-1627483262769-04d0a1401487?q=80&w=100&auto=format&fit=crop' } },
    { id: '3', url: 'https://images.unsplash.com/photo-1531123897727-8f129e1688ce?q=80&w=800&auto=format&fit=crop', height: 250, type: 'photo', user: { name: 'Tigist', avatar: 'https://images.unsplash.com/photo-1531123897727-8f129e1688ce?q=80&w=100&auto=format&fit=crop' } },
    { id: '4', url: 'https://images.unsplash.com/photo-1589156280159-27698a70f29e?q=80&w=800&auto=format&fit=crop', height: 400, type: 'photo', user: { name: 'Hana', avatar: 'https://images.unsplash.com/photo-1589156280159-27698a70f29e?q=80&w=100&auto=format&fit=crop' } },
    { id: '5', url: 'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?q=80&w=800&auto=format&fit=crop', height: 350, type: 'photo', user: { name: 'Abel', avatar: 'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?q=80&w=100&auto=format&fit=crop' } },
    { id: '6', url: 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?q=80&w=800&auto=format&fit=crop', height: 280, type: 'photo', user: { name: 'Marta', avatar: 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?q=80&w=100&auto=format&fit=crop' } },
    { id: '7', url: 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?q=80&w=800&auto=format&fit=crop', height: 500, type: 'photo', user: { name: 'Yodit', avatar: 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?q=80&w=100&auto=format&fit=crop' } },
    { id: '8', url: 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?q=80&w=800&auto=format&fit=crop', height: 320, type: 'photo', user: { name: 'Beti', avatar: 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?q=80&w=100&auto=format&fit=crop' } },
    { id: '9', url: 'https://images.unsplash.com/photo-1529626455594-4ff0802cfb7e?q=80&w=800&auto=format&fit=crop', height: 380, type: 'photo', user: { name: 'Meron', avatar: 'https://images.unsplash.com/photo-1529626455594-4ff0802cfb7e?q=80&w=100&auto=format&fit=crop' } },
    { id: '10', url: 'https://images.unsplash.com/photo-1488426862026-3ee34a7d66df?q=80&w=800&auto=format&fit=crop', height: 290, type: 'photo', user: { name: 'Sara', avatar: 'https://images.unsplash.com/photo-1488426862026-3ee34a7d66df?q=80&w=100&auto=format&fit=crop' } },
    { id: '11', url: 'https://images.unsplash.com/photo-1517841905240-472988babdf9?q=80&w=800&auto=format&fit=crop', height: 420, type: 'photo', user: { name: 'Liya', avatar: 'https://images.unsplash.com/photo-1517841905240-472988babdf9?q=80&w=100&auto=format&fit=crop' } },
    { id: '12', url: 'https://images.unsplash.com/photo-1539571696357-5a69c17a67c6?q=80&w=800&auto=format&fit=crop', height: 310, type: 'photo', user: { name: 'Yonas', avatar: 'https://images.unsplash.com/photo-1539571696357-5a69c17a67c6?q=80&w=100&auto=format&fit=crop' } },
    { id: '13', url: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?q=80&w=800&auto=format&fit=crop', height: 360, type: 'photo', user: { name: 'Kidus', avatar: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?q=80&w=100&auto=format&fit=crop' } },
    { id: '14', url: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?q=80&w=800&auto=format&fit=crop', height: 270, type: 'photo', user: { name: 'Nati', avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?q=80&w=100&auto=format&fit=crop' } },
    { id: '15', url: 'https://images.unsplash.com/photo-1492562080023-ab3db95bfbce?q=80&w=800&auto=format&fit=crop', height: 480, type: 'photo', user: { name: 'Eden', avatar: 'https://images.unsplash.com/photo-1492562080023-ab3db95bfbce?q=80&w=100&auto=format&fit=crop' } },
    { id: '16', url: 'https://images.unsplash.com/photo-1508214751196-bcfd4ca60f91?q=80&w=800&auto=format&fit=crop', height: 330, type: 'photo', user: { name: 'Rahel', avatar: 'https://images.unsplash.com/photo-1508214751196-bcfd4ca60f91?q=80&w=100&auto=format&fit=crop' } },
    { id: '17', url: 'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?q=80&w=800&auto=format&fit=crop', height: 340, type: 'photo', user: { name: 'Daniel', avatar: 'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?q=80&w=100&auto=format&fit=crop' } },
    { id: '18', url: 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?q=80&w=800&auto=format&fit=crop', height: 390, type: 'photo', user: { name: 'Mimi', avatar: 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?q=80&w=100&auto=format&fit=crop' } },
];

interface FeedItem {
    id: string;
    url: string;
    height: number;
    type: 'photo' | 'video';
    user: {
        name: string;
        avatar: string;
    };
    likes?: number;
    comments?: number;
}

export const ExploreFeedScreen = ({ navigation }: any) => {
    const [feed, setFeed] = useState<FeedItem[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [refreshing, setRefreshing] = useState(false);
    const [columns, setColumns] = useState<FeedItem[][]>([[], []]);

    useEffect(() => {
        loadFeed();
    }, []);

    useEffect(() => {
        organizeIntoColumns();
    }, [feed]);

    const loadFeed = async () => {
        try {
            setIsLoading(true);
            const response = await DiscoveryService.getExploreFeed(1, 50);
            if (response.items && response.items.length > 0) {
                const formattedItems = response.items.map((item: any) => {
                    // Calculate height based on media type or use default
                    let itemHeight = 300;
                    if (item.media.media_type === 'video') {
                        itemHeight = 400; // Videos are typically taller
                    } else {
                        // For photos, use a random height between 250-500 for variety
                        itemHeight = 250 + Math.random() * 250;
                    }

                    return {
                        id: item.media.id,
                        url: item.media.url || item.media.thumbnail_url || '',
                        height: itemHeight,
                        type: item.media.media_type,
                        user: {
                            name: item.user.name,
                            avatar: item.user.avatar || '',
                        },
                    };
                });
                setFeed(formattedItems);
            } else {
                // Fallback to mock if no data
                setFeed(MOCK_FEED);
            }
        } catch (error: any) {
            console.error('Load feed error:', error);
            // Use mock data on error for development
            setFeed(MOCK_FEED);
        } finally {
            setIsLoading(false);
        }
    };

    const organizeIntoColumns = () => {
        const cols: FeedItem[][] = [[], []];
        const heights = [0, 0];

        feed.forEach((item) => {
            const shorterColumn = heights[0] <= heights[1] ? 0 : 1;
            cols[shorterColumn].push(item);
            heights[shorterColumn] += item.height;
        });

        setColumns(cols);
    };

    const onRefresh = async () => {
        setRefreshing(true);
        await loadFeed();
        setRefreshing(false);
    };

    const renderItem = (item: FeedItem, columnIndex: number, itemIndex: number) => {
        const itemHeight = item.height;
        const aspectRatio = itemHeight / ITEM_WIDTH;

        return (
            <View key={item.id}>
                <TouchableOpacity
                    style={[styles.item, { width: ITEM_WIDTH, height: ITEM_WIDTH * aspectRatio }]}
                    onPress={() => {
                        // Navigate to full screen view
                        navigation.navigate('ExploreDetail', { item });
                    }}
                    activeOpacity={0.9}
                >
                    <Image source={{ uri: item.url }} style={styles.itemImage} resizeMode="cover" />
                    
                    {/* Overlay gradient */}
                    <LinearGradient
                        colors={['transparent', 'rgba(0,0,0,0.7)']}
                        style={styles.itemOverlay}
                    >
                        {/* User info */}
                        <View style={styles.itemUserInfo}>
                            <Image source={{ uri: item.user.avatar }} style={styles.itemAvatar} />
                            <Text style={styles.itemUserName}>{item.user.name}</Text>
                        </View>

                        {/* Video indicator */}
                        {item.type === 'video' && (
                            <View style={styles.videoIndicator}>
                                <Text style={styles.videoIcon}>▶</Text>
                            </View>
                        )}
                    </LinearGradient>
                </TouchableOpacity>
                {/* Bottom divider */}
                {itemIndex < columns[columnIndex].length - 1 && (
                    <View style={styles.itemDivider} />
                )}
            </View>
        );
    };

    if (isLoading && feed.length === 0) {
        return (
            <SafeAreaView style={styles.container} edges={['top']}>
                <View style={styles.loadingContainer}>
                    <ActivityIndicator size="large" color={COLORS.primary} />
                </View>
            </SafeAreaView>
        );
    }

    return (
        <SafeAreaView style={styles.container} edges={['top']}>
            {/* Header */}
            <View style={styles.header}>
                <Text style={styles.headerTitle}>Explore 🔥</Text>
                <TouchableOpacity
                    style={styles.addVibeButton}
                    onPress={() => navigation.navigate('AddVibe')}
                >
                    <LinearGradient
                        colors={[COLORS.primary, COLORS.primaryDark]}
                        style={styles.addVibeGradient}
                    >
                        <Text style={styles.addVibeIcon}>+</Text>
                        <Text style={styles.addVibeText}>Add Your Vibe</Text>
                    </LinearGradient>
                </TouchableOpacity>
            </View>

            {/* Masonry Grid */}
            <ScrollView
                contentContainerStyle={styles.listContent}
                showsVerticalScrollIndicator={false}
                refreshControl={
                    <RefreshControl
                        refreshing={refreshing}
                        onRefresh={onRefresh}
                        tintColor={COLORS.primary}
                    />
                }
            >
                <View style={styles.masonryContainer}>
                    {columns.map((column, columnIndex) => (
                        <React.Fragment key={columnIndex}>
                            <View style={styles.column}>
                                {column.map((item, itemIndex) => renderItem(item, columnIndex, itemIndex))}
                            </View>
                            {/* Vertical divider between columns */}
                            {columnIndex < columns.length - 1 && (
                                <View style={styles.columnDivider} />
                            )}
                        </React.Fragment>
                    ))}
                </View>
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
    header: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
        padding: SPACING.m,
        borderBottomWidth: 1,
        borderBottomColor: COLORS.surfaceHighlight,
        backgroundColor: COLORS.surface,
    },
    headerTitle: {
        fontSize: 24,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
    },
    addVibeButton: {
        borderRadius: SIZES.radiusM,
        overflow: 'hidden',
    },
    addVibeGradient: {
        flexDirection: 'row',
        alignItems: 'center',
        paddingHorizontal: SPACING.m,
        paddingVertical: SPACING.s,
        gap: SPACING.xs,
    },
    addVibeIcon: {
        fontSize: 18,
        fontWeight: 'bold',
        color: COLORS.background,
    },
    addVibeText: {
        fontSize: 14,
        fontWeight: 'bold',
        color: COLORS.background,
    },
    listContent: {
        padding: 0,
    },
    masonryContainer: {
        flexDirection: 'row',
        alignItems: 'flex-start',
        backgroundColor: COLORS.background,
    },
    column: {
        flex: 1,
    },
    columnDivider: {
        width: DIVIDER_WIDTH,
        backgroundColor: COLORS.background,
    },
    item: {
        borderRadius: 0,
        overflow: 'hidden',
        backgroundColor: COLORS.surface,
    },
    itemDivider: {
        height: DIVIDER_WIDTH,
        backgroundColor: COLORS.background,
        width: '100%',
    },
    itemImage: {
        width: '100%',
        height: '100%',
    },
    itemOverlay: {
        position: 'absolute',
        bottom: 0,
        left: 0,
        right: 0,
        padding: SPACING.s,
    },
    itemUserInfo: {
        flexDirection: 'row',
        alignItems: 'center',
    },
    itemAvatar: {
        width: 24,
        height: 24,
        borderRadius: 12,
        marginRight: SPACING.xs,
        borderWidth: 1,
        borderColor: COLORS.background,
    },
    itemUserName: {
        fontSize: 12,
        fontWeight: '600',
        color: COLORS.background,
    },
    videoIndicator: {
        position: 'absolute',
        top: SPACING.s,
        right: SPACING.s,
        backgroundColor: 'rgba(0,0,0,0.6)',
        borderRadius: 20,
        width: 40,
        height: 40,
        alignItems: 'center',
        justifyContent: 'center',
    },
    videoIcon: {
        fontSize: 16,
        color: COLORS.background,
    },
});

