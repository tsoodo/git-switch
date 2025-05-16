class Gs < Formula
  desc "Git profile switcher for managing multiple GitHub accounts"
  homepage "https://github.com/tsoodo/git-switch"
  version "1.0.0"
  
  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/tsoodo/git-switch/releases/download/v1.0.0/gs-darwin-arm64"
      sha256 "fdd9e65b57552526f83ce5c4b11ef68d0102301599165bcdeb7aed14048c23bc"
    else
      url "https://github.com/tsoodo/git-switch/releases/download/v1.0.0/gs-darwin-amd64"
      sha256 "6901b20126cecdfd27e2f4df5a7e781dad163b818ee95fb30e894fac1da7bbf4"
    end
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://github.com/tsoodo/git-switch/releases/download/v1.0.0/gs-linux-arm64"
      sha256 "3fe9e335008de0118ed98159b953bbf8804d44f8a6e12998c8f2de57515a5177"
    else
      url "https://github.com/tsoodo/git-switch/releases/download/v1.0.0/gs-linux-amd64"
      sha256 "0ab10959e31c165563caca83417f50933a6dcc64039aa60d815362fcd9ddd194"
    end
  end

  def install
    bin.install Dir["gs-*"].first => "gs"
  end

  test do
    system "#{bin}/gs", "help"
  end
end
