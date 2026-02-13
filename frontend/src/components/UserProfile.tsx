import React from 'react';

type User = {
    id: number;
    email: string;
    name: string;
    role: string;
    subscribed_channel: string;
    subscribe_province_id: number;
    subscribe_location_id: number;
};

export default function UserProfile({ user }: { user: User }) {
    return (
        <div className="max-w-md mx-auto bg-white p-6 rounded-lg shadow">
            <div className="flex items-center gap-4 mb-4">
                <div className="w-12 h-12 bg-gray-100 rounded-full flex items-center justify-center">
                    <svg className="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                    </svg>
                </div>
                <div>
                    <div className="text-lg font-semibold">{user.name}</div>
                    <div className="text-sm text-gray-500">{user.email}</div>
                </div>
            </div>

            <ul className="text-sm text-gray-700 space-y-2">
                <li><strong>ID:</strong> {user.id}</li>
                <li><strong>Role:</strong> {user.role}</li>
                <li><strong>Subscribed Channel:</strong> {user.subscribed_channel}</li>
                <li><strong>Subscribe Province ID:</strong> {user.subscribe_province_id}</li>
                <li><strong>Subscribe Location ID:</strong> {user.subscribe_location_id}</li>
            </ul>
        </div>
    );
}
