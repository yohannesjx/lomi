import React, { useState } from 'react';
import {
    View,
    Text,
    StyleSheet,
    TouchableOpacity,
    Image,
    Alert,
    ActivityIndicator,
    ScrollView,
    Platform,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import * as ImagePicker from 'expo-image-picker';
import * as FileSystem from 'expo-file-system';
import { COLORS, SPACING, SIZES } from '../../theme/colors';
import { UserService } from '../../api/services';
import { LinearGradient } from 'expo-linear-gradient';

export const AddVibeScreen = ({ navigation }: any) => {
    const [selectedMedia, setSelectedMedia] = useState<string | null>(null);
    const [mediaType, setMediaType] = useState<'photo' | 'video' | null>(null);
    const [isUploading, setIsUploading] = useState(false);

    const pickImage = async () => {
        const { status } = await ImagePicker.requestMediaLibraryPermissionsAsync();
        if (status !== 'granted') {
            Alert.alert('Permission needed', 'Please grant camera roll permissions');
            return;
        }

        const result = await ImagePicker.launchImageLibraryAsync({
            mediaTypes: ImagePicker.MediaTypeOptions.Images,
            allowsEditing: true,
            aspect: [1, 1],
            quality: 0.8,
        });

        if (!result.canceled && result.assets[0]) {
            setSelectedMedia(result.assets[0].uri);
            setMediaType('photo');
        }
    };

    const pickVideo = async () => {
        const { status } = await ImagePicker.requestMediaLibraryPermissionsAsync();
        if (status !== 'granted') {
            Alert.alert('Permission needed', 'Please grant camera roll permissions');
            return;
        }

        const result = await ImagePicker.launchImageLibraryAsync({
            mediaTypes: ImagePicker.MediaTypeOptions.Videos,
            allowsEditing: true,
            quality: 0.8,
        });

        if (!result.canceled && result.assets[0]) {
            setSelectedMedia(result.assets[0].uri);
            setMediaType('video');
        }
    };

    const takePhoto = async () => {
        const { status } = await ImagePicker.requestCameraPermissionsAsync();
        if (status !== 'granted') {
            Alert.alert('Permission needed', 'Please grant camera permissions');
            return;
        }

        const result = await ImagePicker.launchCameraAsync({
            allowsEditing: true,
            aspect: [1, 1],
            quality: 0.8,
        });

        if (!result.canceled && result.assets[0]) {
            setSelectedMedia(result.assets[0].uri);
            setMediaType('photo');
        }
    };

    const handleUpload = async () => {
        if (!selectedMedia || !mediaType) return;

        setIsUploading(true);
        try {
            // Step 1: Get presigned URL from backend
            const uploadUrlResponse = await UserService.getPresignedUploadURL(mediaType);
            const { upload_url, file_key, headers } = uploadUrlResponse;

            // Step 2: Read file data
            let fileData: string | Blob;
            let contentType: string;

            if (Platform.OS === 'web') {
                // Web: Use fetch to get blob
                const response = await fetch(selectedMedia);
                fileData = await response.blob();
                contentType = mediaType === 'photo' ? 'image/jpeg' : 'video/mp4';
            } else {
                // Mobile: Use FileSystem to read as base64
                const base64 = await FileSystem.readAsStringAsync(selectedMedia, {
                    encoding: FileSystem.EncodingType.Base64,
                });
                fileData = base64;
                contentType = headers?.['Content-Type'] || (mediaType === 'photo' ? 'image/jpeg' : 'video/mp4');
            }

            // Step 3: Upload file to S3/R2 using presigned URL
            let uploadResponse: Response;
            
            if (Platform.OS === 'web') {
                uploadResponse = await fetch(upload_url, {
                    method: 'PUT',
                    body: fileData as Blob,
                    headers: {
                        'Content-Type': contentType,
                    },
                });
            } else {
                // Mobile: Convert base64 to binary
                const binaryString = atob(fileData as string);
                const bytes = new Uint8Array(binaryString.length);
                for (let i = 0; i < binaryString.length; i++) {
                    bytes[i] = binaryString.charCodeAt(i);
                }
                
                uploadResponse = await fetch(upload_url, {
                    method: 'PUT',
                    body: bytes,
                    headers: {
                        'Content-Type': contentType,
                    },
                });
            }

            if (!uploadResponse.ok) {
                throw new Error(`Upload failed: ${uploadResponse.status} ${uploadResponse.statusText}`);
            }

            // Step 4: Create media record in database
            await UserService.uploadMedia({
                media_type: mediaType,
                file_key: file_key,
                display_order: 1,
            });

            Alert.alert(
                'Success! 🎉',
                'Your vibe has been submitted and will appear in the explore feed after approval.',
                [
                    {
                        text: 'OK',
                        onPress: () => {
                            navigation.goBack();
                        },
                    },
                ]
            );
        } catch (error: any) {
            console.error('Upload error:', error);
            Alert.alert('Error', error.response?.data?.error || error.message || 'Failed to upload');
        } finally {
            setIsUploading(false);
        }
    };

    return (
        <SafeAreaView style={styles.container} edges={['top']}>
            <View style={styles.header}>
                <TouchableOpacity onPress={() => navigation.goBack()} style={styles.backButton}>
                    <Text style={styles.backIcon}>←</Text>
                </TouchableOpacity>
                <Text style={styles.headerTitle}>Add Your Vibe ✨</Text>
                <View style={styles.placeholder} />
            </View>

            <ScrollView contentContainerStyle={styles.content} showsVerticalScrollIndicator={false}>
                {selectedMedia ? (
                    <>
                        <View style={styles.previewContainer}>
                            {mediaType === 'photo' ? (
                                <Image source={{ uri: selectedMedia }} style={styles.preview} />
                            ) : (
                                <View style={styles.preview}>
                                    <Text style={styles.videoPlaceholder}>🎥 Video Selected</Text>
                                </View>
                            )}
                        </View>

                        <View style={styles.actions}>
                            <TouchableOpacity
                                style={styles.changeButton}
                                onPress={() => {
                                    setSelectedMedia(null);
                                    setMediaType(null);
                                }}
                            >
                                <Text style={styles.changeButtonText}>Change</Text>
                            </TouchableOpacity>
                            <TouchableOpacity
                                style={[styles.uploadButton, isUploading && styles.uploadButtonDisabled]}
                                onPress={handleUpload}
                                disabled={isUploading}
                            >
                                {isUploading ? (
                                    <ActivityIndicator color={COLORS.background} />
                                ) : (
                                    <Text style={styles.uploadButtonText}>Upload</Text>
                                )}
                            </TouchableOpacity>
                        </View>
                    </>
                ) : (
                    <>
                        <Text style={styles.instructionText}>
                            Share your vibe with the Lomi community! 📸
                        </Text>

                        <View style={styles.optionsContainer}>
                            <TouchableOpacity style={styles.optionCard} onPress={takePhoto}>
                                <LinearGradient
                                    colors={[COLORS.primary, COLORS.primaryDark]}
                                    style={styles.optionGradient}
                                >
                                    <Text style={styles.optionIcon}>📷</Text>
                                    <Text style={styles.optionTitle}>Take Photo</Text>
                                    <Text style={styles.optionDesc}>Capture the moment</Text>
                                </LinearGradient>
                            </TouchableOpacity>

                            <TouchableOpacity style={styles.optionCard} onPress={pickImage}>
                                <View style={styles.optionContent}>
                                    <Text style={styles.optionIcon}>🖼️</Text>
                                    <Text style={styles.optionTitle}>Choose Photo</Text>
                                    <Text style={styles.optionDesc}>From your gallery</Text>
                                </View>
                            </TouchableOpacity>

                            <TouchableOpacity style={styles.optionCard} onPress={pickVideo}>
                                <View style={styles.optionContent}>
                                    <Text style={styles.optionIcon}>🎥</Text>
                                    <Text style={styles.optionTitle}>Choose Video</Text>
                                    <Text style={styles.optionDesc}>Short video clips</Text>
                                </View>
                            </TouchableOpacity>
                        </View>

                        <View style={styles.infoBox}>
                            <Text style={styles.infoTitle}>📋 Guidelines</Text>
                            <Text style={styles.infoText}>
                                • Keep it fun and authentic{'\n'}
                                • No inappropriate content{'\n'}
                                • Photos/videos will be reviewed{'\n'}
                                • Best vibes appear in explore feed
                            </Text>
                        </View>
                    </>
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
    content: {
        padding: SPACING.l,
    },
    instructionText: {
        fontSize: 18,
        color: COLORS.textPrimary,
        textAlign: 'center',
        marginBottom: SPACING.xl,
        lineHeight: 26,
    },
    optionsContainer: {
        gap: SPACING.m,
        marginBottom: SPACING.xl,
    },
    optionCard: {
        borderRadius: SIZES.radiusM,
        overflow: 'hidden',
        marginBottom: SPACING.m,
    },
    optionGradient: {
        padding: SPACING.xl,
        alignItems: 'center',
    },
    optionContent: {
        backgroundColor: COLORS.surface,
        padding: SPACING.xl,
        alignItems: 'center',
        borderWidth: 1,
        borderColor: COLORS.surfaceHighlight,
        borderRadius: SIZES.radiusM,
    },
    optionIcon: {
        fontSize: 48,
        marginBottom: SPACING.m,
    },
    optionTitle: {
        fontSize: 20,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
        marginBottom: SPACING.xs,
    },
    optionDesc: {
        fontSize: 14,
        color: COLORS.textSecondary,
    },
    infoBox: {
        backgroundColor: COLORS.surface,
        borderRadius: SIZES.radiusM,
        padding: SPACING.l,
        marginTop: SPACING.m,
    },
    infoTitle: {
        fontSize: 18,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
        marginBottom: SPACING.m,
    },
    infoText: {
        fontSize: 14,
        color: COLORS.textSecondary,
        lineHeight: 24,
    },
    previewContainer: {
        marginBottom: SPACING.xl,
    },
    preview: {
        width: '100%',
        height: 400,
        borderRadius: SIZES.radiusM,
        backgroundColor: COLORS.surface,
        justifyContent: 'center',
        alignItems: 'center',
    },
    videoPlaceholder: {
        fontSize: 24,
        color: COLORS.textSecondary,
    },
    actions: {
        flexDirection: 'row',
        gap: SPACING.m,
    },
    changeButton: {
        flex: 1,
        backgroundColor: COLORS.surface,
        padding: SPACING.m,
        borderRadius: SIZES.radiusM,
        alignItems: 'center',
        borderWidth: 1,
        borderColor: COLORS.surfaceHighlight,
    },
    changeButtonText: {
        color: COLORS.textPrimary,
        fontSize: 16,
        fontWeight: '600',
    },
    uploadButton: {
        flex: 1,
        backgroundColor: COLORS.primary,
        padding: SPACING.m,
        borderRadius: SIZES.radiusM,
        alignItems: 'center',
    },
    uploadButtonDisabled: {
        opacity: 0.6,
    },
    uploadButtonText: {
        color: COLORS.background,
        fontSize: 16,
        fontWeight: 'bold',
    },
});

