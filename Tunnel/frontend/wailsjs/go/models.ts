export namespace main {
	
	export class SegmentFileAddress {
	    fileSegmentHash: string;
	    segContentSize: number;
	    segFileSize: number;
	    SegmentNumber: number;
	    segAddress: string[];
	
	    static createFrom(source: any = {}) {
	        return new SegmentFileAddress(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fileSegmentHash = source["fileSegmentHash"];
	        this.segContentSize = source["segContentSize"];
	        this.segFileSize = source["segFileSize"];
	        this.SegmentNumber = source["SegmentNumber"];
	        this.segAddress = source["segAddress"];
	    }
	}
	export class TunnelTracerContent {
	    fileHash: string;
	    fileName: string;
	    fileImage: string;
	    fileDescription: string;
	    fileSize: number;
	    allSegmentCount: number;
	    fileExt: string;
	    fileSegments: SegmentFileAddress[];
	
	    static createFrom(source: any = {}) {
	        return new TunnelTracerContent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fileHash = source["fileHash"];
	        this.fileName = source["fileName"];
	        this.fileImage = source["fileImage"];
	        this.fileDescription = source["fileDescription"];
	        this.fileSize = source["fileSize"];
	        this.allSegmentCount = source["allSegmentCount"];
	        this.fileExt = source["fileExt"];
	        this.fileSegments = this.convertValues(source["fileSegments"], SegmentFileAddress);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

