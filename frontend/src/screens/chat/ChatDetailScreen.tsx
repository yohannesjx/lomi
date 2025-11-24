import React, { useState, useEffect, useRef, useCallback } from 'react';
import { View, Text, StyleSheet, FlatList, TextInput, TouchableOpacity, KeyboardAvoidingView, Platform, Image, ActivityIndicator, Alert } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { COLORS, SPACING, SIZES } from '../../theme/colors';
import { ChatService, CoinService } from '../../api/services';
import { GiftService } from '../../api/services';
import { wsService, WSMessage } from '../../api/websocket';
import { useAuthStore } from '../../store/authStore';
import { GiftModal } from '../../components/chat/GiftModal';

interface Message {
    id: string;
    content?: string;
    message_type: 'text' | 'photo' | 'video' | 'voice' | 'gift';
    media_url?: string;
    gift_id?: string;
    sender_id: string;
    receiver_id: string;
    is_read: boolean;
    created_at: string;
    delivery_status?: 'sent' | 'delivered' | 'read';
}

export const ChatDetailScreen = ({ route, navigation }: any) => {
    const { user, match_id } = route.params;
    const { accessToken, user: currentUser } = useAuthStore();
    const [messages, setMessages] = useState<Message[]>([]);
    const [inputText, setInputText] = useState('');
    const [isGiftModalVisible, setIsGiftModalVisible] = useState(false);
    const [isTyping, setIsTyping] = useState(false);
    const [isLoading, setIsLoading] = useState(true);
    const [isSending, setIsSending] = useState(false);
    const [coinBalance, setCoinBalance] = useState(0);
    const typingTimeoutRef = useRef<NodeJS.Timeout | null>(null);
    const flatListRef = useRef<FlatList>(null);

    // Load messages from API
    useEffect(() => {
        if (!match_id) {
            console.error('No match_id provided');
            setIsLoading(false);
            return;
        }

        loadMessages();
        loadCoinBalance();
        connectWebSocket();

        return () => {
            // Cleanup: Send read receipt and disconnect WebSocket listeners
            if (match_id) {
                wsService.sendReadReceipt(match_id);
            }
        };
    }, [match_id]);

    const loadCoinBalance = async () => {
        try {
            const response = await CoinService.getBalance();
            setCoinBalance(response.coin_balance || 0);
        } catch (error: any) {
            console.error('Load coin balance error:', error);
        }
    };

    const loadMessages = async () => {
        try {
            setIsLoading(true);
            const response = await ChatService.getMessages(match_id, 1, 50);
            // Reverse messages to show newest at bottom
            const reversedMessages = response.messages.reverse();
            setMessages(reversedMessages.map((msg: any) => ({
                ...msg,
                delivery_status: msg.is_read ? 'read' : 'delivered',
            })));
        } catch (error: any) {
            console.error('Load messages error:', error);
        } finally {
            setIsLoading(false);
        }
    };

    const connectWebSocket = () => {
        if (!accessToken) {
            console.error('No access token for WebSocket');
            return;
        }

        // Connect WebSocket if not already connected
        wsService.connect(accessToken);

        // Listen for new messages
        const unsubscribeMessage = wsService.on('message', (wsMsg: WSMessage) => {
            if (wsMsg.match_id === match_id) {
                handleIncomingMessage(wsMsg);
            }
        });

        // Listen for typing indicators
        const unsubscribeTyping = wsService.on('typing', (wsMsg: WSMessage) => {
            if (wsMsg.match_id === match_id && wsMsg.sender_id !== currentUser?.id) {
                setIsTyping(wsMsg.is_typing || false);
            }
        });

        // Listen for delivery status
        const unsubscribeDelivery = wsService.on('delivery_status', (wsMsg: WSMessage) => {
            if (wsMsg.match_id === match_id) {
                updateMessageDeliveryStatus(wsMsg.message_id || '', wsMsg.delivery_status || 'sent');
            }
        });

        // Listen for read receipts
        const unsubscribeRead = wsService.on('read_receipt', (wsMsg: WSMessage) => {
            if (wsMsg.match_id === match_id) {
                updateMessageDeliveryStatus(wsMsg.message_id || '', 'read');
            }
        });

        // Send read receipt for existing unread messages
        if (match_id) {
            wsService.sendReadReceipt(match_id);
        }

        return () => {
            unsubscribeMessage();
            unsubscribeTyping();
            unsubscribeDelivery();
            unsubscribeRead();
        };
    };

    const handleIncomingMessage = (wsMsg: WSMessage) => {
        const newMessage: Message = {
            id: wsMsg.message_id || Date.now().toString(),
            content: wsMsg.content as string,
            message_type: (wsMsg.message_type || 'text') as any,
            media_url: wsMsg.media_url,
            gift_id: wsMsg.gift_id,
            sender_id: wsMsg.sender_id || '',
            receiver_id: wsMsg.receiver_id || '',
            is_read: false,
            created_at: wsMsg.timestamp || new Date().toISOString(),
            delivery_status: 'delivered',
        };

        setMessages(prev => [...prev, newMessage]);
        
        // Auto-scroll to bottom
        setTimeout(() => {
            flatListRef.current?.scrollToEnd({ animated: true });
        }, 100);

        // Send read receipt
        if (match_id) {
            wsService.sendReadReceipt(match_id);
        }
    };

    const updateMessageDeliveryStatus = (messageId: string, status: 'sent' | 'delivered' | 'read') => {
        setMessages(prev => prev.map(msg => 
            msg.id === messageId ? { ...msg, delivery_status: status, is_read: status === 'read' } : msg
        ));
    };

    const handleSend = async () => {
        if (!inputText.trim() || isSending || !match_id) return;

        const messageContent = inputText.trim();
        setInputText('');
        setIsSending(true);

        // Optimistically add message
        const tempId = Date.now().toString();
        const optimisticMessage: Message = {
            id: tempId,
            content: messageContent,
            message_type: 'text',
            sender_id: currentUser?.id || '',
            receiver_id: user.id || '',
            is_read: false,
            created_at: new Date().toISOString(),
            delivery_status: 'sent',
        };

        setMessages(prev => [...prev, optimisticMessage]);

        try {
            // Send via WebSocket
            wsService.sendMessage(match_id, messageContent, 'text');

            // Also send via REST API as fallback
            await ChatService.sendMessage({
                match_id: match_id,
                message_type: 'text',
                content: messageContent,
            });

            // Update message with real ID when received from WebSocket
            // The WebSocket will send back the message with the real ID
        } catch (error: any) {
            console.error('Send message error:', error);
            // Remove optimistic message on error
            setMessages(prev => prev.filter(msg => msg.id !== tempId));
        } finally {
            setIsSending(false);
            // Auto-scroll to bottom
            setTimeout(() => {
                flatListRef.current?.scrollToEnd({ animated: true });
            }, 100);
        }
    };

    const handleSendGift = async (gift: any) => {
        if (!match_id || isSending) return;

        setIsSending(true);
        setIsGiftModalVisible(false);

        try {
            // Send gift via new luxury gift API
            const response = await GiftService.sendGiftLuxury({
                receiver_id: user.id,
                gift_type: gift.type || gift.id, // Use type for new system, id for legacy
                match_id: match_id,
            });

            // Update coin balance after sending gift
            await loadCoinBalance();

            // Send gift message via WebSocket
            wsService.sendMessage(match_id, '', 'gift', undefined, gift.type || gift.id);
        } catch (error: any) {
            console.error('Send gift error:', error);
            Alert.alert('Error', error.response?.data?.error || 'Failed to send gift');
        } finally {
            setIsSending(false);
        }
    };

    const handleTyping = (text: string) => {
        setInputText(text);

        // Send typing indicator
        if (match_id && text.length > 0) {
            wsService.sendTyping(match_id, true);

            // Clear previous timeout
            if (typingTimeoutRef.current) {
                clearTimeout(typingTimeoutRef.current);
            }

            // Stop typing after 3 seconds of no input
            typingTimeoutRef.current = setTimeout(() => {
                wsService.sendTyping(match_id, false);
            }, 3000);
        }
    };

    const renderMessage = ({ item }: { item: Message }) => {
        const isMe = item.sender_id === currentUser?.id;
        const time = new Date(item.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });

        return (
            <View style={[
                styles.messageBubble,
                isMe ? styles.myMessage : styles.theirMessage,
                item.message_type === 'gift' && styles.giftMessageBubble
            ]}>
                {item.message_type === 'gift' ? (
                    <Text style={[
                        styles.messageText,
                        isMe ? styles.myMessageText : styles.theirMessageText,
                        styles.giftMessageText
                    ]}>
                        🎁 Sent a gift
                    </Text>
                ) : item.media_url ? (
                    <Image source={{ uri: item.media_url }} style={styles.messageMedia} />
                ) : (
                    <Text style={[
                        styles.messageText,
                        isMe ? styles.myMessageText : styles.theirMessageText
                    ]}>
                        {item.content}
                    </Text>
                )}
                <View style={styles.messageFooter}>
                    <Text style={[
                        styles.messageTime,
                        isMe ? styles.myMessageTime : styles.theirMessageTime
                    ]}>
                        {time}
                    </Text>
                    {isMe && item.delivery_status && (
                        <Text style={styles.deliveryStatus}>
                            {item.delivery_status === 'read' ? '✓✓' : item.delivery_status === 'delivered' ? '✓✓' : '✓'}
                        </Text>
                    )}
                </View>
            </View>
        );
    };

    return (
        <SafeAreaView style={styles.container} edges={['top']}>
            {/* Header */}
            <View style={styles.header}>
                <TouchableOpacity onPress={() => navigation.goBack()} style={styles.backButton}>
                    <Text style={styles.backIcon}>←</Text>
                </TouchableOpacity>

                <View style={styles.headerInfo}>
                    <Image 
                        source={{ uri: user.photo || user.avatar || 'https://via.placeholder.com/150' }} 
                        style={styles.headerAvatar} 
                    />
                    <View>
                        <Text style={styles.headerName}>{user.name || 'Unknown'}</Text>
                        <Text style={styles.headerStatus}>{user.is_online ? 'Online' : 'Offline'}</Text>
                    </View>
                </View>

                <TouchableOpacity style={styles.moreButton}>
                    <Text style={styles.moreIcon}>⋮</Text>
                </TouchableOpacity>
            </View>

            {/* Messages List */}
            {isLoading ? (
                <View style={styles.loadingContainer}>
                    <ActivityIndicator size="large" color={COLORS.primary} />
                </View>
            ) : (
                <FlatList
                    ref={flatListRef}
                    data={messages}
                    renderItem={renderMessage}
                    keyExtractor={item => item.id}
                    contentContainerStyle={styles.messagesList}
                    onContentSizeChange={() => flatListRef.current?.scrollToEnd({ animated: true })}
                    ListFooterComponent={
                        isTyping ? (
                            <View style={[styles.messageBubble, styles.theirMessage]}>
                                <Text style={styles.typingIndicator}>Typing...</Text>
                            </View>
                        ) : null
                    }
                />
            )}

            {/* Input Bar */}
            <KeyboardAvoidingView
                behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
                keyboardVerticalOffset={Platform.OS === 'ios' ? 90 : 0}
            >
                <View style={styles.inputContainer}>
                    <TouchableOpacity
                        style={styles.giftButton}
                        onPress={() => setIsGiftModalVisible(true)}
                    >
                        <Text style={styles.giftIcon}>🎁</Text>
                    </TouchableOpacity>

                    <TextInput
                        style={styles.input}
                        placeholder="Type a message..."
                        placeholderTextColor={COLORS.textSecondary}
                        value={inputText}
                        onChangeText={handleTyping}
                        multiline
                        editable={!isSending}
                    />

                    <TouchableOpacity
                        style={[styles.sendButton, (!inputText.trim() || isSending) && styles.sendButtonDisabled]}
                        onPress={handleSend}
                        disabled={!inputText.trim() || isSending}
                    >
                        {isSending ? (
                            <ActivityIndicator size="small" color={COLORS.background} />
                        ) : (
                            <Text style={styles.sendIcon}>➤</Text>
                        )}
                    </TouchableOpacity>
                </View>
            </KeyboardAvoidingView>

            <GiftModal
                visible={isGiftModalVisible}
                onClose={() => setIsGiftModalVisible(false)}
                onSendGift={handleSendGift}
                coinBalance={coinBalance}
                onBuyCoins={() => navigation.navigate('BuyCoins')}
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
        padding: SPACING.m,
        borderBottomWidth: 1,
        borderBottomColor: COLORS.surfaceHighlight,
        backgroundColor: COLORS.surface,
    },
    backButton: {
        padding: SPACING.s,
        marginRight: SPACING.s,
    },
    backIcon: {
        fontSize: 24,
        color: COLORS.textPrimary,
    },
    headerInfo: {
        flex: 1,
        flexDirection: 'row',
        alignItems: 'center',
    },
    headerAvatar: {
        width: 40,
        height: 40,
        borderRadius: 20,
        marginRight: SPACING.s,
    },
    headerName: {
        fontSize: 16,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
    },
    headerStatus: {
        fontSize: 12,
        color: COLORS.primary,
    },
    moreButton: {
        padding: SPACING.s,
    },
    moreIcon: {
        fontSize: 24,
        color: COLORS.textPrimary,
    },
    messagesList: {
        padding: SPACING.m,
        gap: SPACING.m,
    },
    messageBubble: {
        maxWidth: '80%',
        padding: SPACING.m,
        borderRadius: SIZES.radiusM,
    },
    myMessage: {
        alignSelf: 'flex-end',
        backgroundColor: COLORS.primary,
        borderBottomRightRadius: 4,
    },
    theirMessage: {
        alignSelf: 'flex-start',
        backgroundColor: COLORS.surface,
        borderBottomLeftRadius: 4,
    },
    messageText: {
        fontSize: 16,
        marginBottom: 4,
    },
    myMessageText: {
        color: COLORS.background,
    },
    theirMessageText: {
        color: COLORS.textPrimary,
    },
    messageTime: {
        fontSize: 10,
        alignSelf: 'flex-end',
    },
    myMessageTime: {
        color: 'rgba(0,0,0,0.6)',
    },
    theirMessageTime: {
        color: COLORS.textSecondary,
    },
    inputContainer: {
        flexDirection: 'row',
        alignItems: 'center',
        padding: SPACING.m,
        backgroundColor: COLORS.surface,
        borderTopWidth: 1,
        borderTopColor: COLORS.surfaceHighlight,
    },
    giftButton: {
        padding: SPACING.s,
        marginRight: SPACING.s,
    },
    giftIcon: {
        fontSize: 24,
    },
    input: {
        flex: 1,
        backgroundColor: COLORS.background,
        borderRadius: SIZES.radiusM,
        paddingHorizontal: SPACING.m,
        paddingVertical: SPACING.s,
        color: COLORS.textPrimary,
        maxHeight: 100,
        marginRight: SPACING.s,
    },
    sendButton: {
        backgroundColor: COLORS.primary,
        width: 40,
        height: 40,
        borderRadius: 20,
        alignItems: 'center',
        justifyContent: 'center',
    },
    sendButtonDisabled: {
        backgroundColor: COLORS.surfaceHighlight,
    },
    sendIcon: {
        color: COLORS.background,
        fontSize: 16,
        fontWeight: 'bold',
    },
    giftMessageBubble: {
        backgroundColor: 'rgba(255, 215, 0, 0.1)',
        borderWidth: 1,
        borderColor: COLORS.gold,
    },
    giftMessageText: {
        color: COLORS.gold,
        fontWeight: 'bold',
        fontStyle: 'italic',
    },
    messageMedia: {
        width: 200,
        height: 200,
        borderRadius: SIZES.radiusM,
        marginBottom: 4,
    },
    messageFooter: {
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'flex-end',
        marginTop: 4,
    },
    deliveryStatus: {
        fontSize: 10,
        marginLeft: 4,
        color: COLORS.textSecondary,
    },
    typingIndicator: {
        fontSize: 14,
        color: COLORS.textSecondary,
        fontStyle: 'italic',
    },
    loadingContainer: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
    },
});
