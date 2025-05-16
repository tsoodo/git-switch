class Gs < Formula
  desc "Git profile switcher for managing multiple GitHub accounts"
  homepage "https://github.com/YOUR_USERNAME/gs"
  
  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/YOUR_USERNAME/gs/releases/download/v1.0.0/gs-darwin-arm64"
      sha256 "HASH_FOR_ARM64_BINARY"
    else
      url "https://github.com/YOUR_USERNAME/gs/releases/download/v1.0.0/gs-darwin-amd64"
      sha256 "HASH_FOR_AMD64_BINARY"
    end
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://github.com/YOUR_USERNAME/gs/releases/download/v1.0.0/gs-linux-arm64"
      sha256 "HASH_FOR_LINUX_ARM64_BINARY"
    else
      url "https://github.com/YOUR_USERNAME/gs/releases/download/v1.0.0/gs-linux-amd64"
      sha256 "HASH_FOR_LINUX_AMD64_BINARY"
    end
  end

  def install
    if OS.mac?
      if Hardware::CPU.arm?
        bin.install "gs-darwin-arm64" => "gs"
      else
        bin.install "gs-darwin-amd64" => "gs"
      end
    elsif OS.linux?
      if Hardware::CPU.arm?
        bin.install "gs-linux-arm64" => "gs"
      else
        bin.install "gs-linux-amd64" => "gs"
      end
    end
  end

  test do
    system "#{bin}/gs", "help"
  end
end
