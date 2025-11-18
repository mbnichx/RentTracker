/*
 * -----------------------------------------------------------
 * Author: Madison Nichols
 * Affiliation: WVU Graduate Student
 * Course: SENG 564
 * -----------------------------------------------------------
 */ 

/**
 * Protected tab layout used after authentication.
 *
 * This file configures the bottom tab navigator and protects the routes by
 * checking authentication via the `AuthContext`. If the user is not
 * authenticated they are redirected to the login route.
 */
import { Redirect, Tabs } from 'expo-router';
import React from 'react';
import { useAuth } from "../../contexts/AuthContext";

import { HapticTab } from '@/components/haptic-tab';
import { IconSymbol } from '@/components/ui/icon-symbol';
import { Colors } from '@/constants/theme';
import { useColorScheme } from '@/hooks/use-color-scheme';


export default function TabLayout() {
  // Use the app color scheme to style the active tab color
  const colorScheme = useColorScheme();
  // Access auth state to guard the protected routes
  const { isAuthenticated } = useAuth();

  // Redirect unauthenticated users to the login screen; this keeps the
  // protected screens unreachable via the router when not signed in.
  if (!isAuthenticated) {
    return <Redirect href="./login" />;
  }
  
  return (
    <Tabs
      screenOptions={{
        tabBarActiveTintColor: Colors[colorScheme ?? 'light'].tint,
        headerShown: false,
        // Use a custom tab button that provides haptic feedback
        tabBarButton: HapticTab,
      }}>
      {/* Define each tab screen with an icon and title. The `IconSymbol`
          component centralizes platform-safe icon rendering. */}
      <Tabs.Screen
        name="dashboard"
        options={{
          title: 'Dashboard',
          tabBarIcon: ({ color }) => <IconSymbol size={28} name="house.fill" color={color} />,
        }}
      />
      <Tabs.Screen
        name="rent"
        options={{
          title: 'Rent',
          tabBarIcon: ({ color }) => <IconSymbol size={28} name="dollarsign.circle.fill" color={color} />,
        }}
      />
      <Tabs.Screen
        name="maintenance"
        options={{
          title: 'Maintenance',
          tabBarIcon: ({ color }) => <IconSymbol size={28} name="wrench.fill" color={color} />,
        }}
      />
      <Tabs.Screen
        name="tenants"
        options={{
          title: 'Tenants',
          tabBarIcon: ({ color }) => <IconSymbol size={28} name="person.3.fill" color={color} />,
        }}
      />
      <Tabs.Screen
        name="settings"
        options={{
          title: 'Settings',
          tabBarIcon: ({ color }) => <IconSymbol size={28} name="gearshape.fill" color={color} />,
        }}
      />

    </Tabs>
  );
}
