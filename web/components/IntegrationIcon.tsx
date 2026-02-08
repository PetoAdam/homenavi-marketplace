import { FaSpotify } from 'react-icons/fa';

const faIconMap: Record<string, React.ComponentType<{ className?: string }>> = {
  spotify: FaSpotify
};

type Props = {
  icon?: string;
  fallback: string;
  className?: string;
  imageClassName?: string;
};

export default function IntegrationIcon({ icon, fallback, className, imageClassName }: Props) {
  if (icon?.startsWith('fa:')) {
    const key = icon.slice(3).toLowerCase();
    const FaIcon = faIconMap[key];
    if (FaIcon) {
      return <FaIcon className={className || 'h-6 w-6 text-white/80'} aria-hidden="true" />;
    }
  }

  if (icon) {
    return <img src={icon} alt="" className={imageClassName || 'h-full w-full object-cover'} />;
  }

  return (
    <div className="h-full w-full text-center text-lg leading-[3rem] text-white/60">
      {fallback.slice(0, 1).toUpperCase()}
    </div>
  );
}
