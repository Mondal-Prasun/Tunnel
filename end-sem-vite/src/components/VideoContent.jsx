function VideoContent({ videoUrl }) {
  return (
    <div>
        <video className='w-full' controls>
          {videoUrl ? (
            <source src={videoUrl} type="video/mp4" />
          ) : (
            <p>Video source is unavailable.</p>
          )}
          Your browser does not support the video tag.
        </video>
      </div>
  );
}

export default VideoContent;