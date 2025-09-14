import React from 'react';
import ReactMarkdown from 'react-markdown';

interface MarkdownViewProps {
  children: string;
}

const MarkdownView: React.FC<MarkdownViewProps> = ({ children }) => {
  return <ReactMarkdown>{children}</ReactMarkdown>;
};

export default MarkdownView;
