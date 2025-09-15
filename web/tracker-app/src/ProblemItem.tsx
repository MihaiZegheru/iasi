import React from 'react';
import { Link } from 'react-router-dom';

export interface ProblemItemProps {
  name: string;
  url: string;
  time: string;
  id: string;
  solved: boolean;
  onToggle: () => void;
}


const ProblemItem: React.FC<ProblemItemProps> = ({ name, url, time, id, solved, onToggle }) => (
  <li className="problem-item">
    {/* Left: Custom Checkbox */}
    <div style={{ flex: '0 0 auto', marginRight: 18, display: 'flex', alignItems: 'center' }}>
      <label className="custom-checkbox">
        <input
          type="checkbox"
          checked={solved}
          onChange={onToggle}
        />
        <span className="checkmark"></span>
      </label>
    </div>
    {/* Right: Clickable problem info */}
    <div
      style={{ flex: 1, display: 'flex', flexDirection: 'column', gap: 2, minWidth: 0 }}
      onClick={() => window.location.href = `/problem/${id}`}
    >
      <span style={{
        fontSize: 17,
        fontWeight: 700,
        color: solved ? '#7a7a7a' : '#1a2a3a',
        marginBottom: 1,
        textOverflow: 'ellipsis',
        overflow: 'hidden',
        whiteSpace: 'nowrap',
        letterSpacing: '0.01em',
        lineHeight: 1.2,
      }}>{name}</span>
      <span style={{ color: '#7abaff', fontSize: 13, fontWeight: 500, marginTop: 1 }}>Added: {time}</span>
    </div>
  </li>
);

export default ProblemItem;
