import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, ScrollView, KeyboardAvoidingView, Platform, Alert, Image, StatusBar, ActivityIndicator } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { Button } from '../../components/ui/Button';
import { Input } from '../../components/ui/Input';
import { BackButton } from '../../components/ui/BackButton';
import { COLORS, SPACING, SIZES } from '../../theme/colors';
import { UserService } from '../../api/services';
import { useAuthStore } from '../../store/authStore';

export const EditProfileScreen = ({ navigation }: any) => {
    const { user, refreshUser } = useAuthStore();
    const [name, setName] = useState('');
    const [age, setAge] = useState('');
    const [gender, setGender] = useState<'male' | 'female' | null>(null);
    const [bio, setBio] = useState('');
    const [city, setCity] = useState('');
    const [isSaving, setIsSaving] = useState(false);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        if (user) {
            setName(user.name || '');
            setAge(user.age ? user.age.toString() : '');
            setGender(user.gender || null);
            setBio(user.bio || '');
            setCity(user.city || '');
            setIsLoading(false);
        }
    }, [user]);

    const handleSave = async () => {
        // Validate inputs
        if (!name || !age || !gender) {
            Alert.alert('Missing Fields', 'Please fill in your name, age, and gender.');
            return;
        }

        const ageNum = parseInt(age, 10);
        if (isNaN(ageNum) || ageNum < 18 || ageNum > 100) {
            Alert.alert('Invalid Age', 'Please enter a valid age between 18 and 100.');
            return;
        }

        setIsSaving(true);

        try {
            await UserService.updateProfile({
                name: name.trim(),
                age: ageNum,
                gender,
                bio: bio.trim(),
                city: city.trim(),
            });

            await refreshUser();
            Alert.alert('Success', 'Profile updated successfully!', [
                { text: 'OK', onPress: () => navigation.goBack() }
            ]);
        } catch (error: any) {
            console.error('Update profile error:', error);
            Alert.alert('Error', 'Failed to update profile. Please try again.');
        } finally {
            setIsSaving(false);
        }
    };

    const GenderOption = ({ type, label, icon }: { type: 'male' | 'female', label: string, icon: string }) => (
        <TouchableOpacity
            style={[
                styles.genderOption,
                gender === type && styles.genderOptionSelected
            ]}
            onPress={() => setGender(type)}
        >
            <Text style={styles.genderIcon}>{icon}</Text>
            <Text style={[
                styles.genderLabel,
                gender === type && styles.genderLabelSelected
            ]}>{label}</Text>
        </TouchableOpacity>
    );

    if (isLoading) {
        return (
            <View style={styles.loadingContainer}>
                <ActivityIndicator size="large" color={COLORS.primary} />
            </View>
        );
    }

    return (
        <View style={styles.container}>
            <StatusBar barStyle="dark-content" backgroundColor={COLORS.background} />
            <SafeAreaView style={styles.safeArea} edges={['bottom']}>
                <KeyboardAvoidingView
                    behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
                    style={{ flex: 1 }}
                >
                    <ScrollView contentContainerStyle={styles.scrollContent}>
                        <View style={styles.header}>
                            <BackButton />
                            <Text style={styles.headerTitle}>Edit Profile</Text>
                            <View style={{ width: 40 }} />
                        </View>

                        <View style={styles.form}>
                            <Input
                                label="Full Name"
                                placeholder="e.g. Abebe Bikila"
                                value={name}
                                onChangeText={setName}
                            />

                            <Input
                                label="Age"
                                placeholder="e.g. 24"
                                keyboardType="number-pad"
                                maxLength={2}
                                value={age}
                                onChangeText={setAge}
                            />

                            <Text style={styles.label}>Gender</Text>
                            <View style={styles.genderContainer}>
                                <GenderOption type="male" label="Male" icon="👨🏾" />
                                <GenderOption type="female" label="Female" icon="👩🏾" />
                            </View>

                            <Input
                                label="City"
                                placeholder="e.g. Addis Ababa"
                                value={city}
                                onChangeText={setCity}
                            />

                            <Input
                                label="Bio"
                                placeholder="Tell us about yourself..."
                                value={bio}
                                onChangeText={setBio}
                                multiline
                                numberOfLines={4}
                                style={{ height: 100, textAlignVertical: 'top' }}
                            />
                        </View>

                        <View style={styles.footer}>
                            <Button
                                title={isSaving ? "Saving..." : "Save Changes"}
                                onPress={handleSave}
                                disabled={isSaving}
                                isLoading={isSaving}
                                size="large"
                            />
                        </View>
                    </ScrollView>
                </KeyboardAvoidingView>
            </SafeAreaView>
        </View>
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
        backgroundColor: COLORS.background,
    },
    safeArea: {
        flex: 1,
    },
    scrollContent: {
        flexGrow: 1,
        padding: SPACING.l,
    },
    header: {
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'space-between',
        marginBottom: SPACING.xl,
    },
    headerTitle: {
        fontSize: 20,
        fontWeight: 'bold',
        color: COLORS.textPrimary,
    },
    form: {
        flex: 1,
    },
    label: {
        color: COLORS.textSecondary,
        marginBottom: SPACING.s,
        fontSize: 14,
        fontWeight: '500',
    },
    genderContainer: {
        flexDirection: 'row',
        gap: SPACING.m,
        marginBottom: SPACING.xl,
    },
    genderOption: {
        flex: 1,
        backgroundColor: COLORS.surface,
        borderRadius: SIZES.radiusM,
        padding: SPACING.m,
        alignItems: 'center',
        borderWidth: 1,
        borderColor: COLORS.surfaceHighlight,
    },
    genderOptionSelected: {
        borderColor: COLORS.primary,
        backgroundColor: 'rgba(167, 255, 131, 0.1)',
    },
    genderIcon: {
        fontSize: 32,
        marginBottom: SPACING.s,
    },
    genderLabel: {
        color: COLORS.textSecondary,
        fontWeight: '600',
    },
    genderLabelSelected: {
        color: COLORS.primary,
    },
    footer: {
        marginTop: SPACING.xl,
        marginBottom: SPACING.xl,
    },
});
