import { lazy } from "react";

const Contents = lazy(() => import("@/components/Contents"));

function Assembler() {
  const contents = [
    {
      fileHash:
        "a4406c90d99c8c3be19fa8cf6c5d34c7bd462d61edb98a00f01324cafa21c052",
      fileName: "blossom.mp4",
      fileSize: 62781236,
      allSegmentCount: 6,
      fileExt: ".mp4",
      fileSegments: [
        {
          fileSegmentHash:
            "5ee442944bd3fbb9fe06d494c0585520c666b51da0e0b95e9a857ae6fb54327a",
          segContentSize: 10464256,
          segFileSize: 10464508,
          SegmentNumber: 0,
          segAddress: ["127.0.0.1:6000"],
        },
        {
          fileSegmentHash:
            "21f4552784245410d946d7fed391bf7e68df9cafbda210fbd03e505b5ffb8a3b",
          segContentSize: 10464256,
          segFileSize: 10464508,
          SegmentNumber: 1,
          segAddress: ["127.0.0.1:6000"],
        },
        {
          fileSegmentHash:
            "118ebfb1bfe904eb6d98fa71ba47f300894ea9e73dac724005c8818e6612b192",
          segContentSize: 10464256,
          segFileSize: 10464508,
          SegmentNumber: 2,
          segAddress: ["127.0.0.1:6000"],
        },
        {
          fileSegmentHash:
            "082657d9ee96d7fa232a7867477ebb8e20af4eaef938e5d78f8bcbd9e2c277f8",
          segContentSize: 10464256,
          segFileSize: 10464508,
          SegmentNumber: 3,
          segAddress: ["127.0.0.1:6000"],
        },
        {
          fileSegmentHash:
            "f7dfea510b7a18aed71835067027967c32e7611ea592fc8b4ab6a39c8ac5e9c8",
          segContentSize: 10464256,
          segFileSize: 10464508,
          SegmentNumber: 4,
          segAddress: ["127.0.0.1:6000"],
        },
        {
          fileSegmentHash:
            "3701cd37e60da0b11853b04c33fc135b67f6e789c5dca4d9598fb86c09e0f077",
          segContentSize: 10459956,
          segFileSize: 10464508,
          SegmentNumber: 5,
          segAddress: ["127.0.0.1:6000"],
        },
      ],
    },
  ];
  return (
    <div className="flex flex-col gap-4 w-full h-full p-4">
      <h1 className="font-semibold text-xl">Contents</h1>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 w-full">
        {contents.map((content, index) => (
          <Contents key={index} item={content} onClick={() => setOpen(true)} />
        ))}
      </div>
    </div>
  );
}

export default Assembler;
