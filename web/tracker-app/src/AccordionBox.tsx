import React, { useState } from 'react';

interface AccordionBoxProps {
  title: string;
  children: React.ReactNode;
  defaultOpen?: boolean;
}

const AccordionBox: React.FC<AccordionBoxProps> = ({ title, children, defaultOpen = false }) => {
  const [open, setOpen] = useState(defaultOpen);
  return (
    <div style={{
      border: '1px solid #ccc',
      borderRadius: 8,
      marginBottom: 10,
      boxShadow: open ? '0 2px 8px #eee' : undefined,
      width: '100%',
      flex: '1 1 100%',
      boxSizing: 'border-box',
      overflow: 'hidden',
      display: 'block',
      overflowX: 'hidden',
    }}>
      <div
        onClick={() => setOpen(o => !o)}
        style={{
          cursor: 'pointer',
          background: '#f0f0f0',
          padding: '10px 16px',
          fontWeight: 600,
          borderRadius: open ? '8px 8px 0 0' : 8,
          userSelect: 'none',
          width: '100%',
          flex: '1 1 100%',
          boxSizing: 'border-box',
          overflow: 'hidden',
          display: 'block',
        }}
      >
        {title} <span style={{ float: 'right', fontWeight: 400 }}>{open ? '▲' : '▼'}</span>
      </div>
      {open && (
        <div style={{
          padding: 16,
          background: '#fff',
          borderRadius: '0 0 8px 8px',
          width: '100%',
          maxWidth: '100%',
          minWidth: 0,
          boxSizing: 'border-box',
          overflowWrap: 'break-word',
          wordBreak: 'break-word',
        }}>
          {children}
        </div>
      )}
    </div>
  );
};

export default AccordionBox;
