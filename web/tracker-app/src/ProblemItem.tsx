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
  <li
    className="problem-item"
    style={{
      display: 'flex',
      alignItems: 'center',
      padding: '10px 18px',
      marginBottom: 10,
      borderRadius: 14,
      background: solved ? 'linear-gradient(90deg, #e0e0e0 0%, #f8f8f8 100%)' : 'linear-gradient(90deg, #f7fbff 0%, #eaf6ff 100%)',
      boxShadow: '0 2px 12px 0 rgba(0,0,0,0.07)',
      border: solved ? '1.5px solid #b0b0b0' : '1.5px solid #cce6ff',
      listStyle: 'none',
      transition: 'background 0.2s, box-shadow 0.2s, border 0.2s',
      cursor: 'pointer',
      position: 'relative',
      overflow: 'hidden',
    }}
    onMouseEnter={e => {
      (e.currentTarget as HTMLElement).style.boxShadow = '0 4px 20px 0 rgba(0,180,255,0.13)';
      (e.currentTarget as HTMLElement).style.background = solved
        ? 'linear-gradient(90deg, #e8e8e8 0%, #f4f4f4 100%)'
        : 'linear-gradient(90deg, #f0f8ff 0%, #e0f4ff 100%)';
    }}
    onMouseLeave={e => {
      (e.currentTarget as HTMLElement).style.boxShadow = '0 2px 12px 0 rgba(0,0,0,0.07)';
      (e.currentTarget as HTMLElement).style.background = solved
        ? 'linear-gradient(90deg, #e0e0e0 0%, #f8f8f8 100%)'
        : 'linear-gradient(90deg, #f7fbff 0%, #eaf6ff 100%)';
    }}
  >
    {/* Left: Checkbox */}
    <div style={{ flex: '0 0 auto', marginRight: 18, display: 'flex', alignItems: 'center' }}>
      <input
        type="checkbox"
        checked={solved}
        onChange={onToggle}
        style={{ width: 22, height: 22, cursor: 'pointer', accentColor: '#00b4ff', borderRadius: 6, border: '1.5px solid #b0e0ff', boxShadow: solved ? '0 0 0 2px #b0b0b0' : '0 0 0 2px #cce6ff' }}
      />
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
