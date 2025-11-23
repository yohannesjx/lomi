import React from 'react';
import { View, Text, StyleSheet, Image, TouchableOpacity, Dimensions, ScrollView } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { COLORS, SPACING, SIZES } from '../../theme/colors';

const { width, height } = Dimensions.get('window');

export const ExploreDetailScreen = ({ route, navigation }: any) => {
    const { item } = route.params;

    return (
        <SafeAreaView style={styles.container} edges={['top']}>
            <View style={styles.header}>
                <TouchableOpacity onPress={() => navigation.goBack()} style={styles.backButton}>
                    <Text style={styles.backIcon}>←</Text>
                </TouchableOpacity>
                <View style={styles.headerUser}>
                    <Image source={{ uri: item.user.avatar }} style={styles.headerAvatar} />
                    <Text style={styles.headerName}>{item.user.name}</Text>
                </View>
                <TouchableOpacity style={styles.moreButton}>
                    <Text style={styles.moreIcon}>⋯</Text>
                </TouchableOpacity>
            </View>

            <ScrollView style={styles.content} showsVerticalScrollIndicator={false}>
                <Image source={{ uri: item.url }} style={styles.fullImage} resizeMode="contain" />
                
                <View style={styles.actions}>
                    <TouchableOpacity style={styles.actionButton}>
                        <Text style={styles.actionIcon}>❤️</Text>
                        <Text style={styles.actionText}>Like</Text>
                    </TouchableOpacity>
                    <TouchableOpacity style={styles.actionButton}>
                        <Text style={styles.actionIcon}>💬</Text>
                        <Text style={styles.actionText}>Comment</Text>
                    </TouchableOpacity>
                    <TouchableOpacity style={styles.actionButton}>
                        <Text style={styles.actionIcon}>📤</Text>
                        <Text style={styles.actionText}>Share</Text>
                    </TouchableOpacity>
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
    headerUser: {
        flexDirection: 'row',
        alignItems: 'center',
        flex: 1,
        justifyContent: 'center',
    },
    headerAvatar: {
        width: 32,
        height: 32,
        borderRadius: 16,
        marginRight: SPACING.xs,
    },
    headerName: {
        fontSize: 16,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
    },
    moreButton: {
        padding: SPACING.s,
    },
    moreIcon: {
        fontSize: 24,
        color: COLORS.textPrimary,
    },
    content: {
        flex: 1,
    },
    fullImage: {
        width: width,
        height: height * 0.7,
    },
    actions: {
        flexDirection: 'row',
        justifyContent: 'space-around',
        padding: SPACING.l,
        borderTopWidth: 1,
        borderTopColor: COLORS.surfaceHighlight,
    },
    actionButton: {
        alignItems: 'center',
    },
    actionIcon: {
        fontSize: 28,
        marginBottom: SPACING.xs,
    },
    actionText: {
        fontSize: 12,
        color: COLORS.textSecondary,
    },
});

