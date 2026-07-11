export function newAppPath() {
  return "/app/new";
}

export function appsPath() {
  return "/apps";
}

export function appPath(appSlug: string) {
  return `/app/${appSlug}`;
}

export function appMenuPath(appSlug: string, menu: string) {
  return `${appPath(appSlug)}/${menu}`;
}

export function filesPath(appSlug: string) {
  return appMenuPath(appSlug, "files");
}

export function folderPath(appSlug: string, folderId: string) {
  return `${appPath(appSlug)}/folder/${folderId}`;
}

export function filePath(appSlug: string, fileId: string) {
  return `${appPath(appSlug)}/file/${fileId}`;
}

export function fileDownloadPath(appSlug: string, fileId: string) {
  return `/resources/${appSlug}/files/${fileId}`;
}

export function apiKeysPath(appSlug: string) {
  return appMenuPath(appSlug, "api-keys");
}

export function webhooksPath(appSlug: string) {
  return appMenuPath(appSlug, "webhooks");
}

export function settingsPath(appSlug: string) {
  return appMenuPath(appSlug, "setting");
}

export function settingsMenuPath(appSlug: string, menu: string) {
  return `${appPath(appSlug)}/setting/${menu}`;
}

export function usageMenuPath(appSlug: string) {
  return `${appPath(appSlug)}/setting/usage`;
}

export function usagesMenuPath(appSlug: string, menu: string) {
  return `${usageMenuPath(appSlug)}/${menu}`;
}

export function logoutPath() {
  return "/logout";
}
