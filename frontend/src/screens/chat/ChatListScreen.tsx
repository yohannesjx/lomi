import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, FlatList, Image, TouchableOpacity, RefreshControl } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { COLORS, SPACING, SIZES } from '../../theme/colors';
import { ChatService } from '../../api/services';
import { EmptyState } from '../../components/ui/EmptyState';
import { ChatSkeleton } from '../../components/ui/SkeletonLoader';

export const ChatListScreen = ({ navigation }: any) => {
    const [newMatches, setNewMatches] = useState<any[]>([]);
    const [messages, setMessages] = useState<any[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [refreshing, setRefreshing] = useState(false);

    useEffect(() => {
        loadChats();
    }, []);

    const loadChats = async () => {
        try {
            setIsLoading(true);
            const chatsData = await ChatService.getChats();
            
            // Transform API response to match UI format
            const transformedChats = chatsData.chats.map((chat: any) => ({
                id: chat.match_id,
                match_id: chat.match_id,
                user: {
                    id: chat.user.id,
                    name: chat.user.name,
                    photo: chat.user.photo || 'https://via.placeholder.com/150',
                },
                lastMessage: chat.last_message?.content || 'No messages yet',
                time: chat.last_message?.created_at ? new Date(chat.last_message.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) : '',
                unread: chat.unread_count || 0,
                isGift: chat.last_message?.message_type === 'gift',
            }));

            // Separate new matches (no messages yet) and existing chats
            const newMatchesList = transformedChats.filter((chat: any) => !chat.lastMessage || chat.lastMessage === 'No messages yet');
            const messagesList = transformedChats.filter((chat: any) => chat.lastMessage && chat.lastMessage !== 'No messages yet');

            setNewMatches(newMatchesList);
            setMessages(messagesList);
        } catch (error: any) {
            console.error('Load chats error:', error);
            // Fallback to empty state on error
            setNewMatches([]);
            setMessages([]);
        } finally {
            setIsLoading(false);
            setRefreshing(false);
        }
    };

    const onRefresh = () => {
        setRefreshing(true);
        loadChats();
    };

    const renderMatch = ({ item }: { item: any }) => (
        <TouchableOpacity
            style={styles.matchItem}
            onPress={() => navigation.navigate('ChatDetail', { user: item.user, match_id: item.match_id })}
        >
            <Image source={{ uri: item.user?.photo || item.photo }} style={styles.matchPhoto} />
            <Text style={styles.matchName}>{item.user?.name || item.name}</Text>
        </TouchableOpacity>
    );

    const renderMessage = ({ item }: { item: any }) => (
        <TouchableOpacity
            style={styles.messageItem}
            onPress={() => navigation.navigate('ChatDetail', { user: item.user, match_id: item.match_id })}
        >
            <Image source={{ uri: item.user.photo }} style={styles.avatar} />
            <View style={styles.messageContent}>
                <View style={styles.messageHeader}>
                    <Text style={styles.userName}>{item.user.name}</Text>
                    <Text style={styles.time}>{item.time}</Text>
                </View>
                <Text
                    style={[
                        styles.lastMessage,
                        item.unread > 0 && styles.unreadMessage,
                        item.isGift && styles.giftMessage
                    ]}
                    numberOfLines={1}
                >
                    {item.lastMessage}
                </Text>
            </View>
            {item.unread > 0 && (
                <View style={styles.unreadBadge}>
                    <Text style={styles.unreadText}>{item.unread}</Text>
                </View>
            )}
        </TouchableOpacity>
    );

    return (
        <SafeAreaView style={styles.container}>
            <View style={styles.header}>
                <Text style={styles.title}>Messages</Text>
            </View>

            {isLoading ? (
                <View style={styles.loadingContainer}>
                    <ChatSkeleton />
                    <ChatSkeleton />
                    <ChatSkeleton />
                </View>
            ) : (
                <>
                    {newMatches.length > 0 && (
                        <View style={styles.section}>
                            <Text style={styles.sectionTitle}>New Matches 💚</Text>
                            <FlatList
                                data={newMatches}
                                renderItem={renderMatch}
                                keyExtractor={item => item.id}
                                horizontal
                                showsHorizontalScrollIndicator={false}
                                contentContainerStyle={styles.matchesList}
                            />
                        </View>
                    )}

                    <View style={styles.messagesSection}>
                        {messages.length === 0 ? (
                            <EmptyState
                                icon="💬"
                                title="No messages yet"
                                message="Start swiping to find matches and begin conversations!"
                            />
                        ) : (
                            <FlatList
                                data={messages}
                                renderItem={renderMessage}
                                keyExtractor={item => item.id}
                                contentContainerStyle={styles.messagesList}
                                refreshControl={
                                    <RefreshControl
                                        refreshing={refreshing}
                                        onRefresh={onRefresh}
                                        tintColor={COLORS.primary}
                                    />
                                }
                            />
                        )}
                    </View>
                </>
            )}
        </SafeAreaView>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: COLORS.background,
    },
    header: {
        padding: SPACING.l,
        paddingBottom: SPACING.m,
    },
    title: {
        fontSize: 32,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
    },
    section: {
        marginBottom: SPACING.l,
    },
    sectionTitle: {
        fontSize: 18,
        fontWeight: 'bold',
        color: COLORS.textSecondary,
        marginLeft: SPACING.l,
        marginBottom: SPACING.m,
    },
    matchesList: {
        paddingHorizontal: SPACING.l,
        gap: SPACING.l,
    },
    matchItem: {
        alignItems: 'center',
    },
    matchPhoto: {
        width: 70,
        height: 70,
        borderRadius: 35,
        borderWidth: 2,
        borderColor: COLORS.primary,
        marginBottom: SPACING.xs,
    },
    matchName: {
        color: COLORS.textPrimary,
        fontSize: 14,
        fontWeight: '500',
    },
    messagesSection: {
        flex: 1,
        backgroundColor: COLORS.surface,
        borderTopLeftRadius: SIZES.radiusXL,
        borderTopRightRadius: SIZES.radiusXL,
        paddingTop: SPACING.m,
    },
    messagesList: {
        padding: SPACING.l,
    },
    messageItem: {
        flexDirection: 'row',
        alignItems: 'center',
        marginBottom: SPACING.l,
    },
    avatar: {
        width: 60,
        height: 60,
        borderRadius: 30,
        marginRight: SPACING.m,
    },
    messageContent: {
        flex: 1,
    },
    messageHeader: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        marginBottom: 4,
    },
    userName: {
        fontSize: 16,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
    },
    time: {
        fontSize: 12,
        color: COLORS.textSecondary,
    },
    lastMessage: {
        fontSize: 14,
        color: COLORS.textSecondary,
    },
    unreadMessage: {
        color: COLORS.textPrimary,
        fontWeight: 'bold',
    },
    giftMessage: {
        color: COLORS.secondary,
        fontStyle: 'italic',
    },
    unreadBadge: {
        backgroundColor: COLORS.primary,
        width: 24,
        height: 24,
        borderRadius: 12,
        alignItems: 'center',
        justifyContent: 'center',
        marginLeft: SPACING.s,
    },
    unreadText: {
        color: COLORS.background,
        fontSize: 12,
        fontWeight: 'bold',
    },
    loadingContainer: {
        flex: 1,
        padding: SPACING.l,
    },
});
