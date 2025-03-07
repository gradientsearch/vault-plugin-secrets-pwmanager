interface HvBundles {
  request_id: string;
  lease_id: string;
  lease_duration: number;
  renewable: boolean;
  data: Data;
  warnings: null;
  mount_type: string;
}

interface Data {
  bundles: HvBundle[];
  // Only present on create bundle payload
  path: string;
}

interface HvBundle {
  created: number;
  path: string;
}
