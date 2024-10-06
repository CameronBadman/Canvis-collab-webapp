// src/components/AccountPage.jsx
import React, { useState } from 'react';

function AccountPage() {
  const [user, setUser] = useState({
    name: 'John Doe',
    email: 'john@example.com',
    preferences: {
      newsletter: true,
      notifications: false
    }
  });

  const handlePreferenceChange = (pref) => {
    setUser(prevUser => ({
      ...prevUser,
      preferences: {
        ...prevUser.preferences,
        [pref]: !prevUser.preferences[pref]
      }
    }));
  };

  return (
    <div className="account-page">
      <h1>Your Account</h1>
      <section className="user-info">
        <h2>Personal Information</h2>
        <p><strong>Name:</strong> {user.name}</p>
        <p><strong>Email:</strong> {user.email}</p>
      </section>
      <section className="preferences">
        <h2>Preferences</h2>
        <label>
          <input 
            type="checkbox" 
            checked={user.preferences.newsletter} 
            onChange={() => handlePreferenceChange('newsletter')}
          />
          Receive newsletter
        </label>
        <label>
          <input 
            type="checkbox" 
            checked={user.preferences.notifications} 
            onChange={() => handlePreferenceChange('notifications')}
          />
          Enable notifications
        </label>
      </section>
      <section className="actions">
        <h2>Account Actions</h2>
        <button>Change Password</button>
        <button>Delete Account</button>
      </section>
    </div>
  );
}

export default AccountPage;